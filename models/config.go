package models

import (
	"os"
	"strconv"
)

type Config struct {
	Port       string
	MaxWorkers int
	CacheTtl   int
}

func LoadConfig() *Config {
	return &Config{
		Port:       getEnv("PORT", "8080"),
		MaxWorkers: getEnvAsInt("MAX_WORKERS", 50),
		CacheTtl:   getEnvAsInt("CACHE_TTL", 1),
	}
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
