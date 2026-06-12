package database

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect opens a single shared GORM postgres connection pool.
// All services call this once at startup; the returned *gorm.DB is reused for
// the lifetime of the process — never call this per-request.
//
// Pool budget (per service):
//   - 25 max open  → 3 services × 25 = 75, well under PG default max_connections=100
//   - 5  max idle  → keeps connections warm without holding them under low load
//   - 5m lifetime  → recycles long-lived connections before PG closes them server-side
//   - 1m idle time → returns idle connections to PG faster between traffic bursts
func Connect(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		TranslateError: true,
		Logger:         logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetConnMaxIdleTime(1 * time.Minute)

	log.Println("database connection established")
	return db, nil
}
