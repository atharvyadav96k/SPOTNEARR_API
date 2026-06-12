package services

import (
	"github.com/atharvyadav96k/SPOTNEARR_API/connections/cache"
	"github.com/atharvyadav96k/spotnearr/pkg/mq"
	"gorm.io/gorm"
)

type Services struct {
	UserService   *UserService
	ClaimService  *ClaimService
	ReviewService *ReviewService
	SocialService *SocialService
}

func Init(db *gorm.DB, cache *cache.Cache, publisher *mq.Publisher) *Services {
	base := NewBaseService(db, cache)
	return &Services{
		UserService:   NewUserService(db, cache),
		ClaimService:  NewClaimService(base),
		ReviewService: NewReviewService(base),
		SocialService: NewSocialService(base, publisher),
	}
}
