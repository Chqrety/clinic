package config

type PaymentConfig struct {
	Enabled       bool   // Deferred in this phase
	Provider      string // e.g. "midtrans", "xendit"
	ServerKey     string
	ClientKey     string
	WebhookSecret string
	IsProduction  bool
}

func LoadPaymentConfig() PaymentConfig {
	return PaymentConfig{
		Enabled:      false, // Explicitly postponed / deferred
		Provider:     getEnv("PAYMENT_GATEWAY_PROVIDER", "midtrans"),
		ServerKey:    getEnv("PAYMENT_SERVER_KEY", ""),
		ClientKey:    getEnv("PAYMENT_CLIENT_KEY", ""),
		IsProduction: getEnv("PAYMENT_IS_PRODUCTION", "false") == "true",
	}
}

