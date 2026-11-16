package config

import (
	"fmt"
	"time"
)

type DatabaseConfig struct {
	Host                   string
	Port                   int
	User                   string
	Password               string
	DBName                 string
	SSLMode                string
	MaxOpenConns           int
	MaxIdleConns           int
	ConnMaxLifetimeMinutes time.Duration
}

func LoadDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Host:                   getEnv("DB_HOST", "localhost"),
		Port:                   getEnvInt("DB_PORT", 5432),
		User:                   getEnv("DB_USER", "postgres"),
		Password:               getEnv("DB_PASSWORD", "postgres"),
		DBName:                 getEnv("DB_NAME", "clinic_accounting"),
		SSLMode:                getEnv("DB_SSLMODE", "disable"),
		MaxOpenConns:           getEnvInt("DB_MAX_OPEN_CONNS", 25),
		MaxIdleConns:           getEnvInt("DB_MAX_IDLE_CONNS", 10),
		ConnMaxLifetimeMinutes: time.Duration(getEnvInt("DB_CONN_MAX_LIFETIME_MINUTES", 30)) * time.Minute,
	}
}

func (c DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.SSLMode,
	)
}

