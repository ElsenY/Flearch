package helper

import (
	"context"
	"fmt"
	"time"
)

func RetryRequest(ctx context.Context, fn func() error, maxRetries int, initialBackoff time.Duration) error {
	var err error
	backoff := initialBackoff
	for i := range maxRetries {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if i > 0 {
			fmt.Println("Retry #", i+1)
		}

		err = fn()
		if err == nil {
			return nil
		}

		fmt.Println("Error: ", err, "Retrying in ", backoff, " seconds...")
		backoff *= 2

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoff):
		}
	}

	fmt.Println("Error: ", err, "Max retries reached")
	return err
}
