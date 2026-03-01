# Flearch - Flight Search Aggregator

A backend service that aggregates flight search results from multiple airline providers (AirAsia, Batik Air, Garuda Indonesia, Lion Air) into a unified response format.

### Prerequisites

- Go 1.24.4+

### Running the Application

1. Create a .env file in the root directory, .env.example is provided and can just follow the value (if not created, app will run normally on default config)
2. then run
```bash
go mod tidy
go run app/cmd/main.go
```
app should be running on designated port declared in .env folder (default port 8080)

### Usage
#### Example: Search for Flights with `curl`

```bash
curl -X POST http://localhost:{{PORT_NUMBER}}/api/flight/search \
  -H "Content-Type: application/json" \
```

The endpoint aggregates data; field names and structure may vary by implementation.
Additionally, each flight data has **"best_price_same_flight"** and **"best_amenities_same_flight"** to give best options with the same flight route.

available request body filter field:

1. origin (filter by specified origin)
2. destination (filter by destination)
3. departure_date (filter by departure date)
4. return_date (since there is no return_date data in the mock json, this filter is going to filter by arrival_date)
5. passengers (filter by number of seats available, whether X passengers is less or equal than available seats)
6. cabin_class (filter by cabin class)
7. min_price (filter by price thats greater than X)
8. max_price (filter by price thats lower than X)
9. stops (filter by number of stops)
10. duration_minutes (filter by total durations in minutes, including layover)
11. airline (filter by airline)
12. arrival_date (filter by arrival date, same as return_date filter)
13. limit (limit how many data can be shown, unlimited by default, need to be used with "page")
14. page (fetch X page, need to have "limit" to be enabled)
15. sort_by (sort by available key below)

available sort key (case sensitive):
1. "price"
2. "duration"
3. "departure_time"
4. "arrival_time"

for descending sort, add "-" prefix

_Mock data is located in internal/provider/flightsearch/apimock_

## Design Decisions

- **Provider Pattern**: Each airline API has its own provider implementation, allowing independent parsing of different response formats into a unified model.
- **Controller/Service/Repo Layered Architecture**: Clean separation between HTTP handling, business logic, and data access.
- **Everything in `internal/`**: All application packages are private to this module, enforced by the Go compiler. No `pkg/` directory since this is a standalone service, not a library.
- **In memory caching system**: In memory caching system to get data faster, the cache system simulate redis get and set with key eviction capability.
- **Retry mechanism with exponential timeoff**: Built-in retry system with exponential timeout implemented for every fetch to providers.
- **Timeout System**: Timeout system that is declared by using the value in .env file to decide maximum time spent a request can have
- **Concurrency System to fetch multiple provider**: Built-in Concurrency system to fetch from multiple providers, ensuring faster fetch operation
- **Model and DTO separation**: Model and DTO Separation to limit what data to be shown and to make it easier to modify before reaching end-user.
- **Go Fiber**: Using GoFiber, built on FastAPI for lighting fast HTTP Engine.
