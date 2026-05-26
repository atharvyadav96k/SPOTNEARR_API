package services

import (
	"context"
	"errors"

	"github.com/atharvyadav96k/SPOTNEARR_API/dtos"
	"github.com/atharvyadav96k/SPOTNEARR_API/models"
	"github.com/atharvyadav96k/SPOTNEARR_API/utils/response"
	"gorm.io/gorm"
)

type BusinessService struct {
	base_service
}

func NewBusinessService(db *gorm.DB) *BusinessService {
	return &BusinessService{
		base_service: NewBaseService(db),
	}
}

func (b *BusinessService) RegisterBusiness(business *dtos.Business, userId uint) response.Res {

	registerBusiness := &models.Business{
		BusinessName: business.Name,
		Email:        &business.Email,
		Phone:        &business.Phone,
		UserID:       userId,
	}

	biz, err := b.RepoBusiness().Create(context.Background(), registerBusiness)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return b.ResponseConflict("This account is already associated with the business")
		}
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return b.ResponseConflict("Invalid user account")
		}
		return b.ResponseBadRequest(err.Error())
	}
	return b.ResponseCreated("Business successfully registered", biz)
}

func (b *BusinessService) GetBusinessById(id uint) response.Res {
	business, err := b.RepoBusiness().GetByID(context.Background(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return b.ResponseNotFound("Business not found")
		}
	}
	return b.ResponseOK("Business info", business)
}

func (b *BusinessService) UpdateBusinessById(id uint, business *models.Business) response.Res {
	b.RepoBusiness()
	business.ID = id
	updateBusiness, err := b.RepoBusiness().Update(context.Background(), business)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return b.ResponseNotFound("business not found")
		}
		return b.ResponseBadRequest("Failed to update the business")
	}
	return b.ResponseOK("Business Updated successfully", updateBusiness)
}
