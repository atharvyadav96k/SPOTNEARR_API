package applayer

import (
	"context"
	"net/http"
	"time"

	"github.com/atharvyadav96k/spotnearr/search-svc/config"
	"github.com/atharvyadav96k/spotnearr/search-svc/connections/database"
	"github.com/atharvyadav96k/spotnearr/search-svc/handlers"
	"github.com/atharvyadav96k/spotnearr/search-svc/services"
	syncsvc "github.com/atharvyadav96k/spotnearr/search-svc/sync"
	"gorm.io/gorm"
)

type application struct {
	searchDB    *gorm.DB
	searchHandler *handlers.SearchHandler
	syncTrigger   http.HandlerFunc
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

	poller := syncsvc.NewPoller(vendorDB, searchDB)

	ctx := context.Background()
	go poller.Run(ctx, 30*time.Second)

	searchSvc := services.NewSearchService(searchDB)

	return application{
		searchDB:    searchDB,
		searchHandler: handlers.NewSearchHandler(searchSvc),
		syncTrigger: func(w http.ResponseWriter, r *http.Request) {
			go poller.ProcessPending(context.Background())
			w.WriteHeader(http.StatusAccepted)
		},
	}
}
