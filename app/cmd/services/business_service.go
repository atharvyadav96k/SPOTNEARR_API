package services

import (
	"context"
	"errors"

	"github.com/atharvyadav96k/SPOTNEARR_API/models"
	"github.com/atharvyadav96k/SPOTNEARR_API/repository"
	"github.com/atharvyadav96k/SPOTNEARR_API/utils/response"
	"gorm.io/gorm"
)

type BusinessService struct {
	base_service
	repo repository.IBusinessesRepository
}

func NewBusinessService(repo repository.IBusinessesRepository) *BusinessService {
	return &BusinessService{repo: repo}
}

func (b *BusinessService) RegisterBusiness(business *models.Business) response.Res {
	business, err := b.repo.Create(context.Background(), business)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return b.ResponseConflict("This account is already associated with the business")
		}
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return b.ResponseConflict("Invalid user account")
		}
		return b.ResponseBadRequest(err.Error())
	}
	return b.ResponseCreated("Business successfully registered", business)
}

func (b *BusinessService) GetBusinessById(id uint) response.Res {
	business, err := b.repo.GetByID(context.Background(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return b.ResponseNotFound("Business not found")
		}
	}
	return b.ResponseOK("Business info", business)
}

func (b *BusinessService) UpdateBusinessById(id uint, business *models.Business) response.Res {
	business.ID = id
	updateBusiness, err := b.repo.Update(context.Background(), business)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			b.ResponseNotFound("business not found")
		}
		return b.ResponseBadRequest("Failed to update the business")
	}
	return b.ResponseOK("Business Updated successfully", updateBusiness)
}
