package applayer

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (a *application) NewMux() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/health", a.searchHandler.Health).Methods(http.MethodGet)

	apiV1 := router.PathPrefix("/api/v1").Subrouter()
	apiV1.HandleFunc("/search", a.searchHandler.Search).Methods(http.MethodGet)

	internal := router.PathPrefix("/internal").Subrouter()
	internal.HandleFunc("/sync", a.syncTrigger).Methods(http.MethodPost)

	return router
}
