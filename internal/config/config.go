package config

import "os"

const (
	defaultHTTPPort = "8070"
)

type Config struct {
	HTTPPort     string
	CacheTTLMins int
	PostgresDSN  string
}

func Load() Config {
	return Config{
		HTTPPort:     getEnv("HTTP_PORT", defaultHTTPPort),
		CacheTTLMins: 5,
		PostgresDSN:  getEnv("POSTGRES_DSN", "postgres://app:app@localhost:5432/app?sslmode=disable"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
