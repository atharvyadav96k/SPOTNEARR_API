package app

import (
	"fmt"
	"os"
	"strings"

	"github.com/atharvyadav96k/SPOTNEARR_API/connections/cache"
	"github.com/atharvyadav96k/SPOTNEARR_API/connections/database"
)

func ValidateEnv() error {
	if strings.TrimSpace(os.Getenv("JWT_SECRET")) == "" {
		return fmt.Errorf("JWT_SECRET environment variable is required")
	}
	return nil
}

func (a *App) InitDb() error {
	var err error
	dbURL := os.Getenv("DATABASE_URL")
	if strings.TrimSpace(dbURL) == "" {
		return fmt.Errorf("DATABASE_URL environment variable is required")
	}
	a.db, err = database.InitDB(dbURL)
	if err != nil {
		return err
	}
	if err := database.AutoMigrate(a.db); err != nil {
		return err
	}
	return nil
}

func (a *App) InitCache() error {
	var err error
	cacheUrl := os.Getenv("CACHE_URL")
	password := os.Getenv("CACHE_PASSWORD")
	if strings.TrimSpace(cacheUrl) == "" {
		cacheUrl = "redis://localhost:6379"
	}
	a.cache, err = cache.InitCache(cacheUrl, password)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) InitCaptcha() error {
	captchaUrl := os.Getenv("CAPTCHA_URL")
	if strings.TrimSpace(captchaUrl) == "" {
		return fmt.Errorf("Failed to get captcha url")
	}
	captchaSecret := os.Getenv("CAPTCHA_SECRET_KEY")
	if strings.TrimSpace(captchaSecret) == "" {
		return fmt.Errorf("failed to get captcha secret")
	}
	return nil
}
