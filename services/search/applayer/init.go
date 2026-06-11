package applayer

import (
	"io"
	"log"
	"net/http"

	"github.com/atharvyadav96k/spotnearr/search-svc/config"
	"github.com/atharvyadav96k/spotnearr/search-svc/connections/database"
	"github.com/atharvyadav96k/spotnearr/search-svc/handlers"
	"github.com/atharvyadav96k/spotnearr/search-svc/services"
	syncsvc "github.com/atharvyadav96k/spotnearr/search-svc/sync"
	"gorm.io/gorm"
)

type application struct {
	searchDB      *gorm.DB
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

	applier := syncsvc.NewApplier(searchDB)
	searchSvc := services.NewSearchService(searchDB)

	return application{
		searchDB:      searchDB,
		searchHandler: handlers.NewSearchHandler(searchSvc),
		syncTrigger: func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}
			defer r.Body.Close()
			if err := applier.Apply(r.Context(), body); err != nil {
				log.Printf("sync: apply: %v", err)
				http.Error(w, "sync failed", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusAccepted)
		},
	}
}
