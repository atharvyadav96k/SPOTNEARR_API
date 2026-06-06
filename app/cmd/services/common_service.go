package services

import (
	"github.com/atharvyadav96k/SPOTNEARR_API/connections/cache"
	"gorm.io/gorm"
)

type Services struct {
	UserService      *UserService
	BusinessService  *BusinessService
	InventoryService *InventoryService
	ProductService   *ProductService
	ClaimService     *ClaimService
	SearchService    *SearchService
	CategoryService  *CategoryService
	ReviewService    *ReviewService
}

func Init(db *gorm.DB, cache *cache.Cache) *Services {
	base := NewBaseService(db, cache)
	return &Services{
		UserService:      NewUserService(db, cache),
		BusinessService:  NewBusinessService(db, cache),
		InventoryService: NewInventoryService(db, cache),
		ProductService:   NewProductService(db, cache),
		ClaimService:     NewClaimService(base),
		SearchService:    NewSearchService(db, cache),
		CategoryService:  NewCategoryService(db, cache),
		ReviewService:    NewReviewService(base),
	}
}
