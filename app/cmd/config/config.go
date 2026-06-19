package config

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type Config struct {
	JWTSecret        string
	DatabaseURL      string
	CacheURL         string
	CachePassword    string
	CaptchaURL       string
	CaptchaSecretKey string
	Port             string
	VendorServiceURL string
	RabbitMQURL      string
}

var C *Config

func Load() error {
	var missing []string

	// required := func(key string) string {
	// 	v := strings.TrimSpace(os.Getenv(key))
	// 	if v == "" {
	// 		missing = append(missing, key)
	// 	}
	// 	return v
	// }
	optional := func(key, defaultVal string) string {
		v := strings.TrimSpace(os.Getenv(key))
		if v == "" {
			return defaultVal
		}
		return v
	}

	C = &Config{
		JWTSecret:        optional("JWT_SECRET", "motherfather"),
		DatabaseURL:      optional("DATABASE_URL", "postgresql://admin:admin123@localhost:5432/spotnearr"),
		CaptchaURL:       optional("CAPTCHA_URL", "https://challenges.cloudflare.com/turnstile/v0/siteverify"),
		CaptchaSecretKey: optional("CAPTCHA_SECRET_KEY", "1x0000000000000000000000000000000AA"),
		CacheURL:         optional("CACHE_URL", "redis://localhost:6379"),
		CachePassword:    optional("CACHE_PASSWORD", ""),
		Port:             optional("PORT", "8080"),
		VendorServiceURL: optional("VENDOR_SERVICE_URL", "http://localhost:8081"),
		RabbitMQURL:      optional("RABBITMQ_URL", "amqp://admin:admin123@localhost:5672/"),
	}

	log.Default().Println(*C)

	if len(missing) > 0 {
		return fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}
	return nil
}
