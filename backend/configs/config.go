package configs

import (
	"time"
)

type Config struct {
	App            AppConfig            `json:"app" mapstructure:"app"`
	Database       DatabaseConfig       `json:"database" mapstructure:"database"`
	Redis          RedisConfig          `json:"redis" mapstructure:"redis"`
	Worker         WorkerConfig         `json:"worker" mapstructure:"worker"`
	Auth           AuthConfig           `json:"auth" mapstructure:"auth"`
	Accounting     AccountingConfig     `json:"accounting" mapstructure:"accounting"`
	PaymentGateway PaymentGatewayConfig `json:"payment_gateway" mapstructure:"payment_gateway"`
}

type AppConfig struct {
	Name                   string        `json:"name" mapstructure:"name"`
	Env                    string        `json:"env" mapstructure:"env"`
	Port                   int           `json:"port" mapstructure:"port"`
	ShutdownTimeoutSeconds time.Duration `json:"shutdown_timeout_seconds" mapstructure:"shutdown_timeout_seconds"`
	LogLevel               string        `json:"log_level" mapstructure:"log_level"`
}

type DatabaseConfig struct {
	Host                   string        `json:"host" mapstructure:"host"`
	Port                   int           `json:"port" mapstructure:"port"`
	User                   string        `json:"user" mapstructure:"user"`
	Password               string        `json:"password" mapstructure:"password"`
	DBName                 string        `json:"dbname" mapstructure:"dbname"`
	SSLMode                string        `json:"sslmode" mapstructure:"sslmode"`
	MaxOpenConns           int           `json:"max_open_conns" mapstructure:"max_open_conns"`
	MaxIdleConns           int           `json:"max_idle_conns" mapstructure:"max_idle_conns"`
	ConnMaxLifetimeMinutes time.Duration `json:"conn_max_lifetime_minutes" mapstructure:"conn_max_lifetime_minutes"`
	ConnMaxIdleTimeMinutes time.Duration `json:"conn_max_idle_time_minutes" mapstructure:"conn_max_idle_time_minutes"`
}

type RedisConfig struct {
	Host     string `json:"host" mapstructure:"host"`
	Port     int    `json:"port" mapstructure:"port"`
	Password string `json:"password" mapstructure:"password"`
	DB       int    `json:"db" mapstructure:"db"`
}

type WorkerConfig struct {
	Concurrency int            `json:"concurrency" mapstructure:"concurrency"`
	Queues      map[string]int `json:"queues" mapstructure:"queues"`
}

type AuthConfig struct {
	JWTSecret                 string `json:"jwt_secret" mapstructure:"jwt_secret"`
	AccessTokenExpiryMinutes  int    `json:"access_token_expiry_minutes" mapstructure:"access_token_expiry_minutes"`
	RefreshTokenExpiryDays    int    `json:"refresh_token_expiry_days" mapstructure:"refresh_token_expiry_days"`
}

type AccountingConfig struct {
	BaseCurrency           string `json:"base_currency" mapstructure:"base_currency"`
	FiscalYearStartMonth   int    `json:"fiscal_year_start_month" mapstructure:"fiscal_year_start_month"`
	AutoPosting            bool   `json:"auto_posting" mapstructure:"auto_posting"`
	AllowNegativeInventory bool   `json:"allow_negative_inventory" mapstructure:"allow_negative_inventory"`
}

// PaymentGatewayConfig is configured but deferred (disabled) in Phase 1
type PaymentGatewayConfig struct {
	Enabled       bool   `json:"enabled" mapstructure:"enabled"`
	Provider      string `json:"provider" mapstructure:"provider"`
	ServerKey     string `json:"server_key" mapstructure:"server_key"`
	ClientKey     string `json:"client_key" mapstructure:"client_key"`
	WebhookSecret string `json:"webhook_secret" mapstructure:"webhook_secret"`
	IsProduction  bool   `json:"is_production" mapstructure:"is_production"`
}

