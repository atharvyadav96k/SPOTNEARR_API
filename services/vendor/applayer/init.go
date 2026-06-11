package applayer

import (
	"context"
	"time"

	"github.com/atharvyadav96k/spotnearr/vendor-svc/config"
	"github.com/atharvyadav96k/spotnearr/vendor-svc/connections/cache"
	"github.com/atharvyadav96k/spotnearr/vendor-svc/connections/database"
	"github.com/atharvyadav96k/spotnearr/vendor-svc/events"
	"github.com/atharvyadav96k/spotnearr/vendor-svc/handlers"
	"github.com/atharvyadav96k/spotnearr/vendor-svc/internal_handlers"
	"github.com/atharvyadav96k/spotnearr/vendor-svc/services"
	"gorm.io/gorm"
)

type application struct {
	db              *gorm.DB
	cache           *cache.Cache
	bizHandler      *handlers.BusinessHandler
	invHandler      *handlers.InventoryHandler
	productHandler  *handlers.ProductHandler
	categoryHandler *handlers.CategoryHandler
	internalHandler *internal_handlers.InternalHandler
}

func Init() application {
	if err := config.Load(); err != nil {
		panic(err)
	}

	db, err := database.InitDB(config.C.DatabaseURL)
	if err != nil {
		panic(err)
	}

	if err := database.AutoMigrate(db); err != nil {
		panic(err)
	}

	c, err := cache.InitCache(config.C.CacheURL, config.C.CachePassword)
	if err != nil {
		panic(err)
	}

	flusher := events.NewFlusher(db, config.C.SearchServiceURL)
	go flusher.Run(context.Background(), 30*time.Second)
	notify := func() { go flusher.Flush(context.Background()) }

	bizSvc := services.NewBusinessService(db, c, config.C.UserServiceURL)
	invSvc := services.NewInventoryService(db, c)
	productSvc := services.NewProductService(db, c)
	categorySvc := services.NewCategoryService(db, c)

	return application{
		db:              db,
		cache:           c,
		bizHandler:      handlers.NewBusinessHandler(bizSvc),
		invHandler:      handlers.NewInventoryHandler(invSvc, notify),
		productHandler:  handlers.NewProductHandler(productSvc, notify),
		categoryHandler: handlers.NewCategoryHandler(categorySvc),
		internalHandler: internal_handlers.NewInternalHandler(db),
	}
}
