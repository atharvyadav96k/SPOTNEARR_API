package applayer

import (
	"context"
	"fmt"

	"github.com/atharvyadav96k/spotnearr/pkg/mq"
	"github.com/atharvyadav96k/spotnearr/search-svc/cache"
	"github.com/atharvyadav96k/spotnearr/search-svc/config"
	"github.com/atharvyadav96k/spotnearr/search-svc/handlers"
	"github.com/atharvyadav96k/spotnearr/search-svc/services"
	syncsvc "github.com/atharvyadav96k/spotnearr/search-svc/sync"
	tsclient "github.com/atharvyadav96k/spotnearr/search-svc/typesense"
	"github.com/redis/go-redis/v9"
)

type application struct {
	searchHandler *handlers.SearchHandler
}

func Init() application {
	if err := config.Load(); err != nil {
		panic(err)
	}

	tsClient := tsclient.New(config.C.TypesenseHost, config.C.TypesensePort, config.C.TypesenseAPIKey)
	if err := tsclient.EnsureCollection(context.Background(), tsClient); err != nil {
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
	tcCache := cache.New(rdb)

	indexer := tsclient.NewIndexer(tsClient)
	applier := syncsvc.NewApplier(indexer, tcCache)

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

	searchSvc := services.NewSearchService(tsClient, tcCache)

	return application{
		searchHandler: handlers.NewSearchHandler(searchSvc),
	}
}
