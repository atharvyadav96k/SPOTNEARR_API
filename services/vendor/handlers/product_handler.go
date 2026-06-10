package handlers

import (
	"net/http"

	"github.com/atharvyadav96k/spotnearr/pkg/tokenizer"
	"github.com/atharvyadav96k/spotnearr/vendor-svc/dtos"
	"github.com/atharvyadav96k/spotnearr/vendor-svc/models"
	"github.com/atharvyadav96k/spotnearr/vendor-svc/services"
)

type ProductHandler struct {
	BaseHandler
	svc *services.ProductService
}

func NewProductHandler(svc *services.ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

func (p *ProductHandler) ProductAdd(w http.ResponseWriter, r *http.Request) {
	bizID := p.ClaimGetBusinessID(r)
	var dto dtos.ProductAddRequest
	if err := parseAndValidateBody(r, &dto); err != nil {
		p.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	searchTokens := tokenizer.MergeTokens(dto.Name, dto.Desc)
	product := models.NewProduct(dto.Name, dto.Price, dto.Desc, dto.Quantity, searchTokens)
	for _, id := range dto.CategoryIDs {
		product.Categories = append(product.Categories, models.Category{ID: id})
	}
	p.Response(w, p.svc.AddNewProduct(bizID, product, dto.StoreIDs))
}

func (p *ProductHandler) ProductUpdate(w http.ResponseWriter, r *http.Request) {
	bizID := p.ClaimGetBusinessID(r)
	productID, err := p.GetProductID(r)
	if err != nil {
		p.ResponseBadRequestWithMessage(w, "Invalid product ID")
		return
	}
	var dto dtos.ProductUpdateRequest
	if err := parseAndValidateBody(r, &dto); err != nil {
		p.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	searchTokens := tokenizer.MergeTokens(dto.Name, dto.Desc)
	product := models.NewProduct(dto.Name, dto.Price, dto.Desc, dto.Quantity, searchTokens)
	product.ID = productID
	p.Response(w, p.svc.UpdateProduct(bizID, product))
}

func (p *ProductHandler) ProductDelete(w http.ResponseWriter, r *http.Request) {
	bizID := p.ClaimGetBusinessID(r)
	productID, err := p.GetProductID(r)
	if err != nil {
		p.ResponseBadRequestWithMessage(w, "Invalid product ID")
		return
	}
	p.Response(w, p.svc.DeleteProduct(bizID, productID))
}

func (p *ProductHandler) ProductGet(w http.ResponseWriter, r *http.Request) {
	bizID := p.ClaimGetBusinessID(r)
	productID, err := p.GetProductID(r)
	if err != nil {
		p.ResponseBadRequestWithMessage(w, "Invalid product ID")
		return
	}
	p.Response(w, p.svc.GetProductByID(bizID, productID))
}

func (p *ProductHandler) ProductList(w http.ResponseWriter, r *http.Request) {
	bizID := p.ClaimGetBusinessID(r)
	p.Response(w, p.svc.GetAllProducts(bizID))
}
