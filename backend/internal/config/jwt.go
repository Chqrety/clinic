package config

import "time"

type JWTConfig struct {
	Secret               string
	AccessExpiryMinutes  time.Duration
	RefreshExpiryDays    time.Duration
}

func LoadJWTConfig() JWTConfig {
	return JWTConfig{
		Secret:              getEnv("JWT_SECRET", "super_secret_jwt_encryption_key_change_me"),
		AccessExpiryMinutes: time.Duration(getEnvInt("JWT_ACCESS_EXPIRY_MINUTES", 60)) * time.Minute,
		RefreshExpiryDays:   time.Duration(getEnvInt("JWT_REFRESH_EXPIRY_DAYS", 7)) * 24 * time.Hour,
	}
}

