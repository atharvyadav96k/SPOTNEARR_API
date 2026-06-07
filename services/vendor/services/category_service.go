package services

import (
	"context"
	"errors"

	"github.com/atharvyadav96k/spotnearr/pkg/httputil"
	"github.com/atharvyadav96k/spotnearr/vendor-svc/connections/cache"
	"github.com/atharvyadav96k/spotnearr/vendor-svc/models"
	"gorm.io/gorm"
)

type CategoryService struct{ baseService }

func NewCategoryService(db *gorm.DB, c *cache.Cache) *CategoryService {
	return &CategoryService{baseService: newBaseService(db, c)}
}

func (c *CategoryService) AddCategory(name, slug string) httputil.Res {
	cat := models.Category{Name: name, Slug: slug}
	created, err := c.RepoCategory().CreateCategory(context.Background(), cat)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return c.ResponseConflict("Category already exists")
		}
		return c.ResponseInternalServer("Failed to create category")
	}
	return c.ResponseCreated("Category created", created)
}

func (c *CategoryService) ListCategories() httputil.Res {
	cats, err := c.RepoCategory().GetAllCategories(context.Background())
	if err != nil {
		return c.ResponseInternalServer("Failed to list categories")
	}
	return c.ResponseOK("Categories", cats)
}
