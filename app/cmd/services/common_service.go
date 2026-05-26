package services

import (
	"gorm.io/gorm"
)

type Services struct {
	UserService      *UserService
	BusinessService  *BusinessService
	InventoryService *InventoryService
}

func Init(db *gorm.DB) *Services {
	return &Services{
		UserService:      NewUserService(db),
		BusinessService:  NewBusinessService(db),
		InventoryService: NewInventoryService(db),
	}
}
