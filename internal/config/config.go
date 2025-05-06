package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"go.uber.org/fx"
)

type Config struct {
	Server   ServerConfig
	Roboflow RoboflowConfig
	Log      LogConfig
}

type ServerConfig struct {
	Port            int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

type RoboflowConfig struct {
	APIKey  string
	APIURL  string
	Version string
}

type LogConfig struct {
	Level      string
	Format     string
	OutputPath string
}

func NewConfig() (*Config, error) {
	// Configurações do servidor
	port, _ := strconv.Atoi(getEnv("PORT", "8080"))
	readTimeout, _ := time.ParseDuration(getEnv("READ_TIMEOUT", "5s"))
	writeTimeout, _ := time.ParseDuration(getEnv("WRITE_TIMEOUT", "10s"))
	shutdownTimeout, _ := time.ParseDuration(getEnv("SHUTDOWN_TIMEOUT", "15s"))

	// Configurações do Roboflow
	roboflowAPIKey := getEnv("ROBOFLOW_API_KEY", "")
	if roboflowAPIKey == "" {
		return nil, fmt.Errorf("ROBOFLOW_API_KEY não configurada")
	}

	// Configurações de log
	logLevel := getEnv("LOG_LEVEL", "info")
	logFormat := getEnv("LOG_FORMAT", "json")
	logOutputPath := getEnv("LOG_OUTPUT_PATH", "")

	return &Config{
		Server: ServerConfig{
			Port:            port,
			ReadTimeout:     readTimeout,
			WriteTimeout:    writeTimeout,
			ShutdownTimeout: shutdownTimeout,
		},
		Roboflow: RoboflowConfig{
			APIKey:  roboflowAPIKey,
			APIURL:  getEnv("ROBOFLOW_API_URL", "https://detect.roboflow.com"),
			Version: getEnv("ROBOFLOW_VERSION", "1"),
		},
		Log: LogConfig{
			Level:      logLevel,
			Format:     logFormat,
			OutputPath: logOutputPath,
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// Provider para fx
func ProvideConfig() fx.Option {
	return fx.Provide(NewConfig)
}
