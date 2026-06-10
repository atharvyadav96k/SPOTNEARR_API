package handlers

import (
	"net/http"

	pkgdtos "github.com/atharvyadav96k/spotnearr/pkg/dtos"
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
	var dto pkgdtos.CategoryAddRequest
	if err := parseAndValidateBody(r, &dto); err != nil {
		c.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	c.Response(w, c.svc.AddCategory(dto.Name, dto.Slug))
}
