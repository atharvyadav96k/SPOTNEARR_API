package handlers

import (
	"net/http"

	"github.com/atharvyadav96k/spotnearr/pkg/tokenizer"
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
		p.ResponseBadRequestWithMessage(w, "At least one category is required")
		return
	}
	desc := p.Desc(r)
	storeIDs := p.StoreIDs(r)
	searchTokens := tokenizer.MergeTokens(name, desc)
	product := models.NewProduct(name, *price, desc, *quantity, searchTokens)
	for _, id := range categoryIDs {
		product.Categories = append(product.Categories, models.Category{ID: id})
	}
	p.Response(w, p.svc.AddNewProduct(bizID, product, storeIDs))
}

func (p *ProductHandler) ProductUpdate(w http.ResponseWriter, r *http.Request) {
	bizID := p.ClaimGetBusinessID(r)
	productID, err := p.GetProductID(r)
	if err != nil {
		p.ResponseBadRequestWithMessage(w, "Invalid product ID")
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
	name := p.Name(r)
	desc := p.Desc(r)
	searchTokens := tokenizer.MergeTokens(name, desc)
	product := models.NewProduct(name, *price, desc, *quantity, searchTokens)
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
