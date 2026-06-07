package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/atharvyadav96k/spotnearr/pkg/httputil"
	"github.com/atharvyadav96k/spotnearr/search-svc/services"
)

type SearchHandler struct {
	svc *services.SearchService
}

func NewSearchHandler(svc *services.SearchService) *SearchHandler {
	return &SearchHandler{svc: svc}
}

func (h *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	lat := services.ParseFloat64Param(r.URL.Query().Get("lat"))
	long := services.ParseFloat64Param(r.URL.Query().Get("long"))

	rangeKm := 10.0
	if rv := services.ParseFloat64Param(r.URL.Query().Get("range")); rv != nil {
		rangeKm = *rv
	}

	result := h.svc.Search(q, lat, long, rangeKm)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(result.StatusCode)
	json.NewEncoder(w).Encode(result)
}

func (h *SearchHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

// respond is a shared helper for non-service responses.
func respond(w http.ResponseWriter, statusCode int, body httputil.Res) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(body)
}
