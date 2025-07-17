// internal/config/config.go
package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	AppPort             string `envconfig:"APP_PORT" default:"8080"`
	PrometheusURL       string `envconfig:"PROMETHEUS_URL" required:"true"`
	DBUser              string `envconfig:"DB_USER" required:"true"`
	DBPassword          string `envconfig:"DB_PASSWORD" required:"true"`
	DBHost              string `envconfig:"DB_HOST" required:"true"`
	DBPort              string `envconfig:"DB_PORT" required:"true"`
	DBName              string `envconfig:"DB_NAME" required:"true"`
	SMTPHost            string `envconfig:"SMTP_HOST" required:"true"`
	SMTPPort            int    `envconfig:"SMTP_PORT" required:"true"`
	SMTPUser            string `envconfig:"SMTP_USER" required:"true"`
	SMTPPassword        string `envconfig:"SMTP_PASSWORD" required:"true"`
	AlertRecipientEmail string `envconfig:"ALERT_RECIPIENT_EMAIL" required:"true"`
	GoogleChatWebhookURL string `envconfig:"GOOGLE_CHAT_WEBHOOK_URL"`
}

func LoadConfig() (*Config, error) {
	// Load .env file first. It's okay if it doesn't exist.
	_ = godotenv.Load()

	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}