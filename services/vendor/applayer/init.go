package applayer

import (
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
	publisher       *events.Publisher
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

	pub := events.NewPublisher(c.RedisClient())

	bizSvc := services.NewBusinessService(db, c, config.C.UserServiceURL)
	invSvc := services.NewInventoryService(db, c)
	productSvc := services.NewProductService(db, c)
	categorySvc := services.NewCategoryService(db, c)

	return application{
		db:              db,
		cache:           c,
		publisher:       pub,
		bizHandler:      handlers.NewBusinessHandler(bizSvc),
		invHandler:      handlers.NewInventoryHandler(invSvc),
		productHandler:  handlers.NewProductHandler(productSvc),
		categoryHandler: handlers.NewCategoryHandler(categorySvc),
		internalHandler: internal_handlers.NewInternalHandler(db),
	}
}
