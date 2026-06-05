package handlers

import (
	"net/http"
	"strings"

	"github.com/atharvyadav96k/SPOTNEARR_API/services"
)

type CategoryHandler struct {
	BaseHandler
}

func NewCategoryHandler(services *services.Services) *CategoryHandler {
	return &CategoryHandler{BaseHandler: *NewBaserHandler(services)}
}

func (c *CategoryHandler) CategoryList(w http.ResponseWriter, r *http.Request) {
	res := c.GetCategoryService().ListCategories()
	c.Response(w, res)
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
	res := c.GetCategoryService().AddCategory(name, slug)
	c.Response(w, res)
}
