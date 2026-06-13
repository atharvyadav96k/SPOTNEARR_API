package handlers

import (
	"encoding/json"
	"net/http"

	searchdb "github.com/Developer-Aadesh/spotnearr-database/search"
	pkgdtos "github.com/atharvyadav96k/spotnearr/pkg/dtos"
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
	q := r.URL.Query()
	dto := pkgdtos.NewSearchQuery(
		q.Get("q"), q.Get("lat"), q.Get("long"), q.Get("range"),
		q.Get("category_ids"), q.Get("min_price"), q.Get("max_price"),
	)
	if err := dto.Validate(); err != nil {
		respond(w, http.StatusBadRequest, httputil.Res{Message: err.Error()})
		return
	}
	filters := searchdb.SearchFilters{
		CategoryIDs: dto.CategoryIDs,
		MinPrice:    dto.MinPrice,
		MaxPrice:    dto.MaxPrice,
	}
	result := h.svc.Search(r.Context(), dto.Q, dto.Lat, dto.Long, dto.RangeKm, filters)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(result.StatusCode)
	json.NewEncoder(w).Encode(result)
}

func (h *SearchHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func respond(w http.ResponseWriter, statusCode int, body httputil.Res) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(body)
}
