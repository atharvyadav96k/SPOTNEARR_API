package applayer

import (
	"log"
	"net/http"
	"os"

	pkgmid "github.com/atharvyadav96k/spotnearr/pkg/middleware"
	"github.com/gorilla/mux"
)

func (a *application) NewMux() *mux.Router {
	router := mux.NewRouter()
	router.Use(pkgmid.RequestLogger(log.New(os.Stdout, "[search-svc] ", log.LstdFlags)))

	router.HandleFunc("/health", a.searchHandler.Health).Methods(http.MethodGet)

	apiV1 := router.PathPrefix("/api/v1").Subrouter()
	apiV1.HandleFunc("/search", a.searchHandler.Search).Methods(http.MethodGet)

	return router
}
