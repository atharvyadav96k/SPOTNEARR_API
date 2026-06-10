package handlers

import (
	"net/http"

	"github.com/atharvyadav96k/spotnearr/pkg/tokenizer"
	pkgdtos "github.com/atharvyadav96k/spotnearr/pkg/dtos"
	vendormodel "github.com/Developer-Aadesh/spotnearr-database/vendordb"
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
	var dto pkgdtos.ProductAddRequest
	if err := parseAndValidateBody(r, &dto); err != nil {
		p.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	searchTokens := tokenizer.MergeTokens(dto.Name, dto.Desc)
	product := vendormodel.NewProduct(dto.Name, dto.Price.Value, dto.Price.Unit, dto.Desc, &dto.Quantity.Value, &dto.Quantity.Unit, searchTokens)
	for _, id := range dto.CategoryIDs {
		product.Categories = append(product.Categories, vendormodel.Category{ID: id})
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
	var dto pkgdtos.ProductUpdateRequest
	if err := parseAndValidateBody(r, &dto); err != nil {
		p.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	searchTokens := tokenizer.MergeTokens(dto.Name, dto.Desc)
	product := vendormodel.NewProduct(dto.Name, dto.Price.Value, dto.Price.Unit, dto.Desc, &dto.Quantity.Value, &dto.Quantity.Unit, searchTokens)
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
