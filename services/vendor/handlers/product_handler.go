package handlers

import (
	"errors"
	"net/http"

	vendorpostgres "github.com/Developer-Aadesh/spotnearr-database/vendordb/postgres"
	pkgdtos "github.com/atharvyadav96k/spotnearr/pkg/dtos"
	"github.com/atharvyadav96k/spotnearr/pkg/httputil"
	"github.com/atharvyadav96k/spotnearr/pkg/tokenizer"
	vendormodel "github.com/Developer-Aadesh/spotnearr-database/vendordb"
	"github.com/atharvyadav96k/spotnearr/vendor-svc/services"
	"gorm.io/gorm"
)

type ProductHandler struct {
	BaseHandler
	svc            *services.ProductService
	invProductRepo *vendorpostgres.InvProductRepository
	notify         func()
}

func NewProductHandler(svc *services.ProductService, invProductRepo *vendorpostgres.InvProductRepository, notify func()) *ProductHandler {
	return &ProductHandler{svc: svc, invProductRepo: invProductRepo, notify: notify}
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
	go p.notify()
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
	go p.notify()
}

func (p *ProductHandler) ProductDelete(w http.ResponseWriter, r *http.Request) {
	bizID := p.ClaimGetBusinessID(r)
	productID, err := p.GetProductID(r)
	if err != nil {
		p.ResponseBadRequestWithMessage(w, "Invalid product ID")
		return
	}
	p.Response(w, p.svc.DeleteProduct(bizID, productID))
	go p.notify()
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

// ProductDetail handles GET /api/v1/products/{invProductId}/detail.
// Public endpoint — no auth required. Returns name, price, store info for a
// specific inventory product; intended for the detail view after a search tap.
func (p *ProductHandler) ProductDetail(w http.ResponseWriter, r *http.Request) {
	id, err := p.GetInvProductID(r)
	if err != nil {
		p.ResponseBadRequestWithMessage(w, "invalid product ID")
		return
	}
	detail, err := p.invProductRepo.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respond(w, http.StatusNotFound, httputil.Res{Message: "not found"})
			return
		}
		respond(w, http.StatusInternalServerError, httputil.Res{Message: "internal error"})
		return
	}
	respond(w, http.StatusOK, httputil.Res{Message: "product detail", Data: detail})
}
