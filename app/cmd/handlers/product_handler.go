package handlers

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/factories/token"
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
		p.ResponseBadRequestWithMessage(w, "Price is required")
		return
	}
	quantity := p.Quantity(r)
	if quantity == nil {
		p.ResponseBadRequestWithMessage(w, "Quantity is required")
		return
	}
	categoryIDs := p.CategoryIDs(r)
	if len(categoryIDs) == 0 {
		p.ResponseBadRequestWithMessage(w, "at least one category is required")
		return
	}
	storeIDs := p.StoreIDs(r)
	desc := p.Desc(r)
	searchToken := token.MergeTokens(name, desc)
	product := models.NewProduct(name, *price, desc, *quantity, searchToken)
	for _, id := range categoryIDs {
		product.Categories = append(product.Categories, models.Category{ID: id})
	}
	res := p.GetProductService().AddNewProduct(bizID, product, storeIDs)
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
		p.ResponseBadRequestWithMessage(w, "Price is required")
		return
	}
	desc := p.Desc(r)
	quantity := p.Quantity(r)
	if quantity == nil {
		p.ResponseBadRequestWithMessage(w, "Quantity is required")
		return
	}
	searchToken := token.MergeTokens(name, desc)
	product := models.NewProduct(name, *price, desc, *quantity, searchToken)
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
