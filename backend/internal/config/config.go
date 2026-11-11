package config

import (
	"os"
	"strconv"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Payment  PaymentConfig
}

type AppConfig struct {
	Name     string
	Env      string
	Port     int
	LogLevel string
}

func LoadConfig() *Config {
	return &Config{
		App: AppConfig{
			Name:     getEnv("APP_NAME", "clinic-accounting"),
			Env:      getEnv("APP_ENV", "development"),
			Port:     getEnvInt("APP_PORT", 8080),
			LogLevel: getEnv("LOG_LEVEL", "info"),
		},
		Database: LoadDatabaseConfig(),
		JWT:      LoadJWTConfig(),
		Payment:  LoadPaymentConfig(),
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return defaultVal
}

