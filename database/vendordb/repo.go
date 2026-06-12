package vendordb

import "context"

type IBusinessesRepository interface {
	Create(ctx context.Context, business *Business) (*Business, error)
	GetByID(ctx context.Context, id uint) (*Business, error)
	Update(ctx context.Context, businessID uint, businessName string, desc string) (*Business, error)
	Delete(ctx context.Context, id uint) error
	HardDelete(ctx context.Context, id uint) error
	GetByEmail(ctx context.Context, email string) (*Business, error)
	GetByPhone(ctx context.Context, phone string) (*Business, error)
	VerifyBusiness(ctx context.Context, id uint) error
	ToggleActiveStatus(ctx context.Context, id uint, isActive bool) error
	IncrementFollowerCount(ctx context.Context, id uint) error
	DecrementFollowerCount(ctx context.Context, id uint) error
}

type IStoreRepository interface {
	GetStoreByBusinessId(ctx context.Context, bizID uint) ([]Store, error)
	CreateStoreByBusinessId(ctx context.Context, store *Store) error
	UpdateStoreByBusinessId(ctx context.Context, store Store) (Store, error)
}

type IProductRepository interface {
	AddProduct(ctx context.Context, product Product, bizID uint) (Product, error)
	UpdateProduct(ctx context.Context, product Product, bizID uint) (Product, error)
	GetProductsByBusiness(ctx context.Context, bizID uint) ([]Product, error)
	GetProductById(ctx context.Context, id uint, bizID uint) (Product, error)
	DeleteProduct(ctx context.Context, productID uint, bizID uint) error
}

// InvProductDetail is the snapshot returned by the internal endpoint for claim validation.
type InvProductDetail struct {
	ID            uint    `json:"id"`
	Available     bool    `json:"available"`
	ProductName   string  `json:"product_name"`
	Price         float64 `json:"price"`
	PriceUnit     string  `json:"price_unit"`
	StoreName     string  `json:"store_name"`
	StreetAddress string  `json:"street_address"`
	Lat           float64 `json:"lat"`
	Long          float64 `json:"long"`
}

type IInventoryProduct interface {
	AddProduct(ctx context.Context, invProduct InventoryProduct, bizID uint) error
	GetInvProductList(ctx context.Context, storeID uint, bizID uint) ([]Product, error)
	UpdateProduct(ctx context.Context, invProduct InventoryProduct, bizID uint) (InventoryProduct, error)
	RemoveProduct(ctx context.Context, invProdID uint, bizID uint) error
	GetByID(ctx context.Context, id uint) (*InvProductDetail, error)
}

type IAccessRepository interface {
	GetAccessByUserId(ctx context.Context, id uint) (BusinessAccess, error)
	GetAccessByBusinessId(ctx context.Context, id uint) (BusinessAccess, error)
	CreateNewAccess(ctx context.Context, access BusinessAccess) error
	RemoveAccessByUserId(ctx context.Context, id uint) error
}

type ICategoryRepository interface {
	CreateCategory(ctx context.Context, category Category) (Category, error)
	GetAllCategories(ctx context.Context) ([]Category, error)
	GetCategoryByIDs(ctx context.Context, ids []uint) ([]Category, error)
}

type IProductTokenRepository interface {
	GetTokenCategoryFreqs(ctx context.Context, tokens []string) (map[uint]int, error)
}
