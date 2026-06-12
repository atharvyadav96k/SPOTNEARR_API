package applayer

import (
	"context"
	"time"

	vendorpostgres "github.com/Developer-Aadesh/spotnearr-database/vendordb/postgres"
	"github.com/atharvyadav96k/spotnearr/pkg/mq"
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
	authHandler     *handlers.AuthHandler
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

	mqConn, err := mq.Connect(config.C.RabbitMQURL)
	if err != nil {
		panic(err)
	}
	mqPub, err := mq.NewPublisher(mqConn)
	if err != nil {
		panic(err)
	}
	mqSub, err := mq.NewSubscriber(mqConn)
	if err != nil {
		panic(err)
	}
	if err := events.InitConsumers(context.Background(), mqSub, db); err != nil {
		panic(err)
	}

	flusher := events.NewFlusher(db, mqPub)
	go flusher.Run(context.Background(), 30*time.Second)
	notify := func() { go flusher.Flush(context.Background()) }

	authSvc := services.NewAuthService(db, config.C.JWTSecret)
	bizSvc := services.NewBusinessService(db, c, config.C.UserServiceURL)
	invSvc := services.NewInventoryService(db, c)
	productSvc := services.NewProductService(db, c)
	categorySvc := services.NewCategoryService(db, c)

	return application{
		db:              db,
		cache:           c,
		authHandler:     handlers.NewAuthHandler(authSvc),
		bizHandler:      handlers.NewBusinessHandler(bizSvc),
		invHandler:      handlers.NewInventoryHandler(invSvc, notify),
		productHandler:  handlers.NewProductHandler(productSvc, vendorpostgres.NewInvProductRepository(db), notify),
		categoryHandler: handlers.NewCategoryHandler(categorySvc),
		internalHandler: internal_handlers.NewInternalHandler(db),
	}
}
