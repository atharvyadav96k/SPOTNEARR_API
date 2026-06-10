package postgres

import (
	"context"

	"github.com/Developer-Aadesh/spotnearr-database/vendordb"
	"gorm.io/gorm"
)

type CategoryRepository struct{ db *gorm.DB }

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (c *CategoryRepository) CreateCategory(ctx context.Context, category vendordb.Category) (vendordb.Category, error) {
	err := c.db.WithContext(ctx).Create(&category).Error
	return category, err
}

func (c *CategoryRepository) GetAllCategories(ctx context.Context) ([]vendordb.Category, error) {
	var cats []vendordb.Category
	err := c.db.WithContext(ctx).Order("name ASC").Find(&cats).Error
	return cats, err
}

func (c *CategoryRepository) GetCategoryByIDs(ctx context.Context, ids []uint) ([]vendordb.Category, error) {
	var cats []vendordb.Category
	err := c.db.WithContext(ctx).Where("id IN ?", ids).Find(&cats).Error
	return cats, err
}
