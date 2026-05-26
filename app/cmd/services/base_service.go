package services

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/repository"
	"github.com/atharvyadav96k/SPOTNEARR_API/repository/implementation"
	"github.com/atharvyadav96k/SPOTNEARR_API/utils/response"
	"gorm.io/gorm"
)

type repo struct {
	userRepo   repository.IUserRepository
	bizRepo    repository.IBusinessesRepository
	storeRepo  repository.IStoreRepository
	accessRepo repository.IAccessRepository
}

type base_service struct {
	repo repo
	db   *gorm.DB
}

func NewBaseService(db *gorm.DB) base_service {
	return base_service{
		db: db,
		repo: repo{
			userRepo:   implementation.NewUserRepository(db),
			bizRepo:    implementation.NewBusinessRepository(db),
			storeRepo:  implementation.NewStoreRepository(db),
			accessRepo: implementation.NewAccessRepository(db),
		},
	}
}

func (b *base_service) RepoBusiness() repository.IBusinessesRepository {
	return b.repo.bizRepo
}

func (b *base_service) RepoUser() repository.IUserRepository {
	return b.repo.userRepo
}

func (b *base_service) RepoStore() repository.IStoreRepository {
	return b.repo.storeRepo
}

func (b *base_service) RepoAccess() repository.IAccessRepository {
	return b.repo.accessRepo
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
