package services

import (
	"net/http"

	"github.com/atharvyadav96k/spotnearr/pkg/httputil"
	"github.com/atharvyadav96k/spotnearr/vendor-svc/connections/cache"
	"github.com/atharvyadav96k/spotnearr/vendor-svc/repository"
	"github.com/atharvyadav96k/spotnearr/vendor-svc/repository/implementation"
	"gorm.io/gorm"
)

type repos struct {
	bizRepo          repository.IBusinessesRepository
	storeRepo        repository.IStoreRepository
	accessRepo       repository.IAccessRepository
	invProductRepo   repository.IInventoryProduct
	productRepo      repository.IProductRepository
	categoryRepo     repository.ICategoryRepository
	productTokenRepo repository.IProductTokenRepository
}

type baseService struct {
	cache *cache.Cache
	repo  repos
}

func newBaseService(db *gorm.DB, c *cache.Cache) baseService {
	return baseService{
		cache: c,
		repo: repos{
			bizRepo:          implementation.NewBusinessRepository(db),
			storeRepo:        implementation.NewStoreRepository(db),
			accessRepo:       implementation.NewAccessRepository(db),
			invProductRepo:   implementation.NewInvProductRepository(db),
			productRepo:      implementation.NewProductRepository(db),
			categoryRepo:     implementation.NewCategoryRepository(db),
			productTokenRepo: implementation.NewProductTokenRepository(db),
		},
	}
}

func (b *baseService) Cache() *cache.Cache        { return b.cache }
func (b *baseService) RepoBusiness() repository.IBusinessesRepository { return b.repo.bizRepo }
func (b *baseService) RepoStore() repository.IStoreRepository         { return b.repo.storeRepo }
func (b *baseService) RepoAccess() repository.IAccessRepository       { return b.repo.accessRepo }
func (b *baseService) RepoInvProduct() repository.IInventoryProduct   { return b.repo.invProductRepo }
func (b *baseService) RepoProduct() repository.IProductRepository     { return b.repo.productRepo }
func (b *baseService) RepoCategory() repository.ICategoryRepository   { return b.repo.categoryRepo }

func res(message string, statusCode int, data interface{}) httputil.Res {
	return httputil.NewResponse(message, statusCode, data)
}

func (b *baseService) ResponseOK(message string, data interface{}) httputil.Res {
	return res(message, http.StatusOK, data)
}
func (b *baseService) ResponseCreated(message string, data interface{}) httputil.Res {
	return res(message, http.StatusCreated, data)
}
func (b *baseService) ResponseNotFound(message string) httputil.Res {
	return res(message, http.StatusNotFound, nil)
}
func (b *baseService) ResponseBadRequest(message string) httputil.Res {
	return res(message, http.StatusBadRequest, nil)
}
func (b *baseService) ResponseConflict(message string) httputil.Res {
	return res(message, http.StatusConflict, nil)
}
func (b *baseService) ResponseInternalServer(message string) httputil.Res {
	return res(message, http.StatusInternalServerError, nil)
}
func (b *baseService) ResponseUnauthorized() httputil.Res {
	return res("unauthorized", http.StatusUnauthorized, nil)
}
