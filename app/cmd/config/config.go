package config

import (
	"fmt"
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
}

var C *Config

func Load() error {
	var missing []string

	required := func(key string) string {
		v := strings.TrimSpace(os.Getenv(key))
		if v == "" {
			missing = append(missing, key)
		}
		return v
	}
	optional := func(key, defaultVal string) string {
		v := strings.TrimSpace(os.Getenv(key))
		if v == "" {
			return defaultVal
		}
		return v
	}

	C = &Config{
		JWTSecret:        required("JWT_SECRET"),
		DatabaseURL:      required("DATABASE_URL"),
		CaptchaURL:       required("CAPTCHA_URL"),
		CaptchaSecretKey: required("CAPTCHA_SECRET_KEY"),
		CacheURL:         optional("CACHE_URL", "redis://localhost:6379"),
		CachePassword:    optional("CACHE_PASSWORD", ""),
		Port:             optional("PORT", "8080"),
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}
	return nil
}
