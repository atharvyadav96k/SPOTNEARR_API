package applayer

import (
	"context"
	"time"

	"github.com/atharvyadav96k/spotnearr/search-svc/config"
	"github.com/atharvyadav96k/spotnearr/search-svc/connections/database"
	"github.com/atharvyadav96k/spotnearr/search-svc/handlers"
	"github.com/atharvyadav96k/spotnearr/search-svc/services"
	syncsvc "github.com/atharvyadav96k/spotnearr/search-svc/sync"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type application struct {
	searchDB      *gorm.DB
	searchHandler *handlers.SearchHandler
}

func Init() application {
	if err := config.Load(); err != nil {
		panic(err)
	}

	searchDB, err := database.InitDB(config.C.DatabaseURL)
	if err != nil {
		panic(err)
	}
	if err := database.AutoMigrate(searchDB); err != nil {
		panic(err)
	}

	vendorDB, err := database.InitDB(config.C.VendorDBURL)
	if err != nil {
		panic(err)
	}

	redisOpt, err := redis.ParseURL(config.C.CacheURL)
	if err != nil {
		panic(err)
	}
	if config.C.CachePassword != "" {
		redisOpt.Password = config.C.CachePassword
	}
	redisClient := redis.NewClient(redisOpt)

	poller := syncsvc.NewPoller(vendorDB, searchDB)
	subscriber := syncsvc.NewSubscriber(redisClient, poller)

	ctx := context.Background()
	go subscriber.Run(ctx)
	go poller.Run(ctx, 30*time.Second)

	searchSvc := services.NewSearchService(searchDB)

	return application{
		searchDB:      searchDB,
		searchHandler: handlers.NewSearchHandler(searchSvc),
	}
}
