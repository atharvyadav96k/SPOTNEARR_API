package services

import (
	"context"
	"errors"

	"github.com/atharvyadav96k/SPOTNEARR_API/connections/cache"
	"github.com/atharvyadav96k/SPOTNEARR_API/models"
	"github.com/atharvyadav96k/SPOTNEARR_API/utils/response"
	"gorm.io/gorm"
)

type ProductService struct {
	base_service
}

func NewProductService(db *gorm.DB, cache *cache.Cache) *ProductService {
	return &ProductService{
		base_service: NewBaseService(db, cache),
	}
}

func (p ProductService) GetProductByID(bizID uint, productID uint) response.Res {
	product, err := p.RepoProduct().GetProductById(context.Background(), productID, bizID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return p.ResponseNotFound("product not found")
		}
		return p.ResponseInternalServer("Failed to get product")
	}
	return p.ResponseOK("Product", product)
}

func (p *ProductService) AddNewProduct(bizID uint, product models.Product) response.Res {
	err := p.RepoProduct().AddProduct(context.Background(), product, bizID)
	if err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return p.ResponseBadRequest("Business does not exists")
		}
		return p.ResponseInternalServer("Failed to add product")
	}
	return p.ResponseCreated("Product added successfully", nil)
}

func (p *ProductService) UpdateProduct(bizID uint, product models.Product) response.Res {
	product, err := p.RepoProduct().UpdateProduct(context.Background(), product, bizID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return p.ResponseBadRequest("Product not found")
		}
		return p.ResponseInternalServer("Failed to update product")
	}
	return p.ResponseOK("Product updated successfully", product)
}

func (p *ProductService) DeleteProduct(bizID uint, productID uint) response.Res {
	err := p.RepoProduct().DeleteProduct(context.Background(), productID, bizID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return p.ResponseNotFound("Product not found")
		}
		return p.ResponseInternalServer("Failed to delete product")
	}
	return p.ResponseOK("Product remove successfully", nil)
}
