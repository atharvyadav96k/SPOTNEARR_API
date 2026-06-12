package applayer

import (
	"context"
	"fmt"
	"time"

	"github.com/atharvyadav96k/spotnearr/pkg/mq"
	"github.com/atharvyadav96k/spotnearr/search-svc/config"
	"github.com/atharvyadav96k/spotnearr/search-svc/connections/database"
	freqpkg "github.com/atharvyadav96k/spotnearr/search-svc/freq"
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

	opts, err := redis.ParseURL(config.C.CacheURL)
	if err != nil {
		panic(fmt.Errorf("search: redis url: %w", err))
	}
	if config.C.CachePassword != "" {
		opts.Password = config.C.CachePassword
	}
	rdb := redis.NewClient(opts)

	flusher := freqpkg.NewFlusher(searchDB, rdb)
	go flusher.Run(context.Background(), 30*time.Second)

	applier := syncsvc.NewApplier(searchDB, rdb)

	mqConn, err := mq.Connect(config.C.RabbitMQURL)
	if err != nil {
		panic(err)
	}
	mqSub, err := mq.NewSubscriber(mqConn)
	if err != nil {
		panic(err)
	}
	if err := syncsvc.RunConsumer(context.Background(), mqSub, applier); err != nil {
		panic(err)
	}

	searchSvc := services.NewSearchService(searchDB, rdb)

	return application{
		searchDB:      searchDB,
		searchHandler: handlers.NewSearchHandler(searchSvc),
	}
}
