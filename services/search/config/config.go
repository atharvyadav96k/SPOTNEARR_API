package config

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type Config struct {
	JWTSecret     string
	DatabaseURL   string
	VendorDBURL   string
	CacheURL      string
	CachePassword string
	Port          string
}

var C *Config

func Load() error {
	optional := func(key, defaultVal string) string {
		v := strings.TrimSpace(os.Getenv(key))
		if v == "" {
			return defaultVal
		}
		return v
	}

	C = &Config{
		JWTSecret:     optional("JWT_SECRET", "motherfather"),
		DatabaseURL:   optional("DATABASE_URL", "postgresql://admin:admin123@localhost:5432/spotnearr_search"),
		VendorDBURL:   optional("VENDOR_DB_URL", "postgresql://admin:admin123@localhost:5432/spotnearr"),
		CacheURL:      optional("CACHE_URL", "redis://localhost:6379"),
		CachePassword: optional("CACHE_PASSWORD", ""),
		Port:          optional("PORT", "8080"),
	}

	log.Default().Println(*C)

	var missing []string
	if len(missing) > 0 {
		return fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}
	return nil
}
