package services

import (
	"github.com/atharvyadav96k/SPOTNEARR_API/connections/cache"
	"gorm.io/gorm"
)

type Services struct {
	UserService      *UserService
	BusinessService  *BusinessService
	InventoryService *InventoryService
}

func Init(db *gorm.DB, cache *cache.Cache) *Services {
	return &Services{
		UserService:      NewUserService(db, cache),
		BusinessService:  NewBusinessService(db, cache),
		InventoryService: NewInventoryService(db, cache),
	}
}
