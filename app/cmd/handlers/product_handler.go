package handlers

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/models"
	"github.com/atharvyadav96k/SPOTNEARR_API/services"
)

type ProductHandler struct {
	BaseHandler
}

func NewProductHandler(service *services.Services) *ProductHandler {
	return &ProductHandler{
		BaseHandler: *NewBaserHandler(service),
	}
}

func (p *ProductHandler) ProductAdd(w http.ResponseWriter, r *http.Request) {
	bizID := p.ClaimGetBusinessId(r)
	name := p.Name(r)
	if name == "" {
		p.ResponseBadRequestWithMessage(w, "Product name is required")
		return
	}
	price := p.Price(r)
	if price == nil {
		p.ResponseBadRequestWithMessage(w, "Product price is required")
		return
	}
	if *price < 1 {
		p.ResponseBadRequestWithMessage(w, "Price should not be 0")
		return
	}
	product := models.NewProduct(name, *price)
	res := p.GetProductService().AddNewProduct(bizID, product)
	p.Response(w, res)
}

func (p *ProductHandler) ProductUpdate(w http.ResponseWriter, r *http.Request) {
	bizID := p.ClaimGetBusinessId(r)
	productId, err := p.GetProductId(r)
	if err != nil {
		p.ResponseBadRequestWithMessage(w, "failed to get product")
		return
	}
	name := p.Name(r)
	price := p.Price(r)
	if price == nil {
		p.ResponseBadRequestWithMessage(w, "Product price is required")
		return
	}
	product := models.NewProduct(name, *price)
	product.ID = productId
	res := p.GetProductService().UpdateProduct(bizID, product)
	p.Response(w, res)
}

func (p *ProductHandler) ProductDelete(w http.ResponseWriter, r *http.Request) {
	bizID := p.ClaimGetBusinessId(r)
	productID, err := p.GetProductId(r)
	if err != nil {
		p.ResponseBadRequestWithMessage(w, "failed to get product")
		return
	}
	res := p.GetProductService().DeleteProduct(bizID, productID)
	p.Response(w, res)
}

func (p *ProductHandler) ProductGet(w http.ResponseWriter, r *http.Request) {
	bizId := p.ClaimGetBusinessId(r)
	productID, err := p.GetProductId(r)
	if err != nil {
		p.ResponseBadRequestWithMessage(w, "failed to get product")
		return
	}
	res := p.GetProductService().GetProductByID(bizId, productID)
	p.Response(w, res)
}

func (p *ProductHandler) ProductOfBusiness(w http.ResponseWriter, r *http.Request) {
	p.ResponseOK(w)
}

func (p *ProductHandler) ProductNearBy(w http.ResponseWriter, r *http.Request) {
	p.ResponseOK(w)
}
