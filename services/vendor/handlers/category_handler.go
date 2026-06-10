package handlers

import (
	"net/http"
	"strings"

	"github.com/atharvyadav96k/spotnearr/vendor-svc/services"
)

type CategoryHandler struct {
	BaseHandler
	svc *services.CategoryService
}

func NewCategoryHandler(svc *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: svc}
}

func (c *CategoryHandler) CategoryList(w http.ResponseWriter, r *http.Request) {
	c.Response(w, c.svc.ListCategories())
}

func (c *CategoryHandler) CategoryAdd(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimSpace(c.Name(r))
	if name == "" {
		c.ResponseBadRequestWithMessage(w, "name is required")
		return
	}
	slug := strings.TrimSpace(c.Slug(r))
	if slug == "" {
		c.ResponseBadRequestWithMessage(w, "slug is required")
		return
	}
	c.Response(w, c.svc.AddCategory(name, slug))
}
