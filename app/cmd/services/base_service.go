package services

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/connections/cache"
	"github.com/atharvyadav96k/SPOTNEARR_API/utils/response"
	usermodel "github.com/Developer-Aadesh/spotnearr-database/user"
	userpostgres "github.com/Developer-Aadesh/spotnearr-database/user/postgres"
	"gorm.io/gorm"
)

type repo struct {
	userRepo   usermodel.IUserRepository
	claimRepo  usermodel.IClaimRepository
	reviewRepo usermodel.IReviewRepository
}

type base_service struct {
	cache *cache.Cache
	repo  repo
	db    *gorm.DB
}

func NewBaseService(db *gorm.DB, cache *cache.Cache) base_service {
	return base_service{
		db:    db,
		cache: cache,
		repo: repo{
			userRepo:   userpostgres.NewUserRepository(db),
			claimRepo:  userpostgres.NewClaimRepository(db),
			reviewRepo: userpostgres.NewReviewRepository(db),
		},
	}
}

func (b *base_service) Cache() *cache.Cache {
	return b.cache
}

func (b *base_service) RepoUser() usermodel.IUserRepository {
	return b.repo.userRepo
}

func (b *base_service) RepoClaim() usermodel.IClaimRepository {
	return b.repo.claimRepo
}

func (b *base_service) RepoReview() usermodel.IReviewRepository {
	return b.repo.reviewRepo
}

func res(message string, statusCode int, data interface{}) response.Res {
	return response.Res{
		Message:    message,
		StatusCode: statusCode,
		Data:       data,
	}
}

func (b *base_service) ResponseOK(message string, data interface{}) response.Res {
	return res(message, http.StatusOK, data)
}

func (b *base_service) ResponseNotFound(message string) response.Res {
	return res(message, http.StatusNotFound, nil)
}

func (b *base_service) ResponseInternalServer(message string) response.Res {
	return res(message, http.StatusInternalServerError, nil)
}

func (b *base_service) ResponseBadRequest(message string) response.Res {
	return res(message, http.StatusBadRequest, nil)
}

func (b *base_service) ResponseUnauthorized() response.Res {
	return res("unauthorized", http.StatusUnauthorized, nil)
}

func (b *base_service) ResponseConflict(message string) response.Res {
	return res(message, http.StatusConflict, nil)
}

func (b *base_service) ResponseCreated(message string, data interface{}) response.Res {
	return res(message, http.StatusCreated, data)
}
