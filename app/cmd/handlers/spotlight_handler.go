package handlers

import "net/http"

type SpotlightHandler struct {
	BaseHandler
}

func NewSpotlightHandler() *SpotlightHandler {
	return &SpotlightHandler{}
}

func (s *SpotlightHandler) SpotlightAdd(w http.ResponseWriter, r *http.Request) {
	s.ResponseOK(w)
}

func (s *SpotlightHandler) SpotlightUpdate(w http.ResponseWriter, r *http.Request) {
	s.ResponseOK(w)
}

func (s *SpotlightHandler) SpotlightDelete(w http.ResponseWriter, r *http.Request) {
	s.ResponseOK(w)
}
