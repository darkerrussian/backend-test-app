package config

import "os"

const (
	defaultHTTPPort = "8070"
)

type Config struct {
	HTTPPort     string
	CacheTTLMins int
}

func Load() Config {
	return Config{
		HTTPPort:     getEnv("HTTP_PORT", defaultHTTPPort),
		CacheTTLMins: 5,
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
