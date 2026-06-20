package config

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type Config struct {
	TypesenseHost   string
	TypesensePort   string
	TypesenseAPIKey string
	CacheURL        string
	CachePassword   string
	Port            string
	RabbitMQURL     string
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
		TypesenseHost:   optional("TYPESENSE_HOST", "localhost"),
		TypesensePort:   optional("TYPESENSE_PORT", "8108"),
		TypesenseAPIKey: optional("TYPESENSE_API_KEY", "dev-api-key"),
		CacheURL:        optional("CACHE_URL", "redis://localhost:6379"),
		CachePassword:   optional("CACHE_PASSWORD", ""),
		Port:            optional("PORT", "8080"),
		RabbitMQURL:     optional("RABBITMQ_URL", "amqp://admin:admin123@localhost:5672/"),
	}

	log.Printf("search config: typesense=%s:%s port=%s", C.TypesenseHost, C.TypesensePort, C.Port)

	var missing []string
	if C.TypesenseAPIKey == "" {
		missing = append(missing, "TYPESENSE_API_KEY")
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}
	return nil
}
