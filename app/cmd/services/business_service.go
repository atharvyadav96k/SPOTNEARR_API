package services

import (
	"context"
	"errors"
	"log"

	"github.com/atharvyadav96k/SPOTNEARR_API/connections/cache"
	"github.com/atharvyadav96k/SPOTNEARR_API/dtos"
	"github.com/atharvyadav96k/SPOTNEARR_API/models"
	"github.com/atharvyadav96k/SPOTNEARR_API/utils/response"
	"gorm.io/gorm"
)

type BusinessService struct {
	base_service
}

func NewBusinessService(db *gorm.DB, cache *cache.Cache) *BusinessService {
	return &BusinessService{
		base_service: NewBaseService(db, cache),
	}
}

func (b *BusinessService) RegisterBusiness(business *dtos.Business, userId uint) response.Res {
	biz, err := b.RepoBusiness().Create(context.Background(), models.NewBusiness(business.Name, business.Email, business.Phone, business.Desc))
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return b.ResponseConflict("Phone number or email is already in use")
		}
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return b.ResponseConflict("Invalid user account")
		}
		return b.ResponseBadRequest(err.Error())
	}

	err = b.RepoAccess().CreateNewAccess(context.Background(), models.NewBusinessAccess(userId, biz.ID, models.RoleAdmin))
	if err != nil {
		log.Default().Println(err)
		b.RepoBusiness().HardDelete(context.Background(), biz.ID)
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return b.ResponseConflict("This account is already associated with the business")
		}
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return b.ResponseConflict("Invalid user account")
		}
		return b.ResponseBadRequest(err.Error())
	}
	if err := b.Cache().GetRefreshTokenSession().InvalidateRefreshToken(userId); err != nil {
		log.Default().Println("Failed to invalidate refresh token after business registration:", err)
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

func (b *BusinessService) UpdateBusinessById(id uint, businessId uint, businessName string, desc string) response.Res {
	updateBusiness, err := b.RepoBusiness().Update(context.Background(), businessId, businessName, desc)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return b.ResponseNotFound("business not found")
		}
		return b.ResponseBadRequest("Failed to update the business")
	}
	return b.ResponseOK("Business Updated successfully", updateBusiness)
}
