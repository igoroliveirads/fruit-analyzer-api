package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port           int
	AllowedOrigins []string
	MaxRequestSize int64
	Environment    string
}

func NewConfig() *Config {
	port, _ := strconv.Atoi(getEnv("PORT", "8080"))
	maxRequestSize, _ := strconv.ParseInt(getEnv("MAX_REQUEST_SIZE", "10485760"), 10, 64) // 10MB default

	return &Config{
		Port:           port,
		AllowedOrigins: []string{getEnv("ALLOWED_ORIGINS", "*")},
		MaxRequestSize: maxRequestSize,
		Environment:    getEnv("ENVIRONMENT", "development"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
