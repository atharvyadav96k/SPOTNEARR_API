package handlers

import "net/http"

type ClaimHandler struct {
	BaseHandler
}

func NewClaimHandler() *ClaimHandler {
	return &ClaimHandler{}
}

func (c *ClaimHandler) ClaimProduct(w http.ResponseWriter, r *http.Request) {
	c.ResponseOK(w)
}

func (c *ClaimHandler) ClaimRemove(w http.ResponseWriter, r *http.Request) {
	c.ResponseOK(w)
}
