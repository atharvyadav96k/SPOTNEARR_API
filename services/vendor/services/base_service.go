package services

import (
	"net/http"

	"github.com/atharvyadav96k/spotnearr/pkg/httputil"
	"github.com/atharvyadav96k/spotnearr/vendor-svc/connections/cache"
	vendormodel "github.com/Developer-Aadesh/spotnearr-database/vendordb"
	vendorpostgres "github.com/Developer-Aadesh/spotnearr-database/vendordb/postgres"
	"gorm.io/gorm"
)

type repos struct {
	bizRepo          vendormodel.IBusinessesRepository
	storeRepo        vendormodel.IStoreRepository
	accessRepo       vendormodel.IAccessRepository
	invProductRepo   vendormodel.IInventoryProduct
	productRepo      vendormodel.IProductRepository
	categoryRepo     vendormodel.ICategoryRepository
	productTokenRepo vendormodel.IProductTokenRepository
}

type baseService struct {
	cache *cache.Cache
	repo  repos
}

func newBaseService(db *gorm.DB, c *cache.Cache) baseService {
	return baseService{
		cache: c,
		repo: repos{
			bizRepo:          vendorpostgres.NewBusinessRepository(db),
			storeRepo:        vendorpostgres.NewStoreRepository(db),
			accessRepo:       vendorpostgres.NewAccessRepository(db),
			invProductRepo:   vendorpostgres.NewInvProductRepository(db),
			productRepo:      vendorpostgres.NewProductRepository(db),
			categoryRepo:     vendorpostgres.NewCategoryRepository(db),
			productTokenRepo: vendorpostgres.NewProductTokenRepository(db),
		},
	}
}

func (b *baseService) Cache() *cache.Cache                                { return b.cache }
func (b *baseService) RepoBusiness() vendormodel.IBusinessesRepository   { return b.repo.bizRepo }
func (b *baseService) RepoStore() vendormodel.IStoreRepository           { return b.repo.storeRepo }
func (b *baseService) RepoAccess() vendormodel.IAccessRepository         { return b.repo.accessRepo }
func (b *baseService) RepoInvProduct() vendormodel.IInventoryProduct     { return b.repo.invProductRepo }
func (b *baseService) RepoProduct() vendormodel.IProductRepository       { return b.repo.productRepo }
func (b *baseService) RepoCategory() vendormodel.ICategoryRepository     { return b.repo.categoryRepo }

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
