package config

// load all configs from environment variables
type Config struct {
	ServerConfig ServerConfig
	CacheConfig  CacheConfig
	AppConfig    AppConfig
}

func NewConfig() *Config {
	return &Config{
		ServerConfig: NewServerConfig(),
		CacheConfig:  NewCacheConfig(),
		AppConfig:    NewAppConfig(),
	}
}
