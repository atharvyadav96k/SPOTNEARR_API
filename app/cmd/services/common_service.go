package services

import (
	"github.com/atharvyadav96k/SPOTNEARR_API/connections/cache"
	"gorm.io/gorm"
)

type Services struct {
	UserService   *UserService
	ClaimService  *ClaimService
	ReviewService *ReviewService
}

func Init(db *gorm.DB, cache *cache.Cache) *Services {
	base := NewBaseService(db, cache)
	return &Services{
		UserService:   NewUserService(db, cache),
		ClaimService:  NewClaimService(base),
		ReviewService: NewReviewService(base),
	}
}
