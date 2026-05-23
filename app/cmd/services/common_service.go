package services

import (
	"github.com/atharvyadav96k/SPOTNEARR_API/repository/implementation"
	"gorm.io/gorm"
)

type Services struct {
	UserService     *UserService
	BusinessService *BusinessService
}

func Init(db *gorm.DB) *Services {
	return &Services{
		UserService:     NewUserService(implementation.NewUserRepository(db)),
		BusinessService: NewBusinessService(implementation.NewBusinessRepository(db)),
	}
}
