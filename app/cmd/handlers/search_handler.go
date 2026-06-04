package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/atharvyadav96k/SPOTNEARR_API/services"
)

type SearchHandler struct {
	BaseHandler
}

func NewSearchHandler(services *services.Services) *SearchHandler {
	return &SearchHandler{
		BaseHandler: *NewBaserHandler(services),
	}
}

func (s *SearchHandler) QuerySearch(r *http.Request) string {
	return r.URL.Query().Get("search")
}

func (s *SearchHandler) queryFloat(r *http.Request, key string) *float64 {
	v := r.URL.Query().Get(key)
	if v == "" {
		return nil
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return nil
	}
	return &f
}

func (s *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	search := strings.TrimSpace(s.QuerySearch(r))
	if search == "" {
		s.ResponseBadRequest(w)
		return
	}
	lat := s.queryFloat(r, "lat")
	long := s.queryFloat(r, "long")
	rangeKm := 5.0
	if v := s.queryFloat(r, "range"); v != nil {
		rangeKm = *v
	}
	res := s.GetSearchService().Search(search, lat, long, rangeKm)
	s.Response(w, res)
}
