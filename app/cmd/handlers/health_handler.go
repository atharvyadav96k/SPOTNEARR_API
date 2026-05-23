package handlers

import "net/http"

type Health struct {
	BaseHandler
}

func NewHealthHandler() *Health {
	return &Health{}
}

func (h *Health) HealthOK(w http.ResponseWriter, r *http.Request) {
	h.ResponseOK(w)
}
