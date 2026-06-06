package implementation

import (
	"context"

	"github.com/atharvyadav96k/SPOTNEARR_API/models"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (c *CategoryRepository) CreateCategory(ctx context.Context, category models.Category) (models.Category, error) {
	err := c.db.WithContext(ctx).Create(&category).Error
	return category, err
}

func (c *CategoryRepository) GetAllCategories(ctx context.Context) ([]models.Category, error) {
	var cats []models.Category
	err := c.db.WithContext(ctx).Order("name ASC").Find(&cats).Error
	return cats, err
}

func (c *CategoryRepository) GetCategoryByIDs(ctx context.Context, ids []uint) ([]models.Category, error) {
	var cats []models.Category
	err := c.db.WithContext(ctx).Where("id IN ?", ids).Find(&cats).Error
	return cats, err
}
