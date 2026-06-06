package services

import (
	"context"
	"errors"

	"github.com/atharvyadav96k/SPOTNEARR_API/connections/cache"
	"github.com/atharvyadav96k/SPOTNEARR_API/models"
	"github.com/atharvyadav96k/SPOTNEARR_API/utils/response"
	"gorm.io/gorm"
)

type CategoryService struct {
	base_service
}

func NewCategoryService(db *gorm.DB, cache *cache.Cache) *CategoryService {
	return &CategoryService{base_service: NewBaseService(db, cache)}
}

func (c *CategoryService) AddCategory(name, slug string) response.Res {
	cat := models.Category{Name: name, Slug: slug}
	created, err := c.RepoCategory().CreateCategory(context.Background(), cat)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return c.ResponseConflict("category already exists")
		}
		return c.ResponseInternalServer("failed to create category")
	}
	return c.ResponseCreated("category created", created)
}

func (c *CategoryService) ListCategories() response.Res {
	cats, err := c.RepoCategory().GetAllCategories(context.Background())
	if err != nil {
		return c.ResponseInternalServer("failed to list categories")
	}
	return c.ResponseOK("categories", cats)
}
