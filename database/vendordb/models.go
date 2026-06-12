package vendordb

import (
	"time"

	"github.com/atharvyadav96k/spotnearr/pkg/jwtutil"
	"gorm.io/gorm"
)

type Business struct {
	ID           uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	BusinessName string `gorm:"type:varchar(255);not null" json:"businessName"`

	Email *string `gorm:"type:varchar(255);unique;index" json:"email"`
	Phone *string `gorm:"type:varchar(20);unique;index" json:"phone"`
	Desc  string  `gorm:"type:varchar(100)" json:"description"`

	IsActive         bool `gorm:"default:true;not null" json:"isActive"`
	VerifiedBusiness bool `gorm:"default:false;not null" json:"verifiedBusiness"`
	FollowerCount    uint `gorm:"default:0;not null" json:"followerCount"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

func NewBusiness(name string, email string, phone string, desc string) *Business {
	return &Business{
		BusinessName: name,
		Email:        &email,
		Phone:        &phone,
		Desc:         desc,
	}
}

// ---------------------------------------------------------------------------

type Store struct {
	ID            uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name          string `gorm:"type:varchar(255);not null" json:"name"`
	StreetAddress string `gorm:"type:text;not null" json:"streetAddress"`

	BusinessID uint      `gorm:"index;not null" json:"businessId"`
	Business   *Business `gorm:"foreignKey:BusinessID;constraint:OnDelete:CASCADE;" json:"business,omitempty"`

	Lat  float64 `gorm:"type:decimal(10,8);not null" json:"lat"`
	Long float64 `gorm:"type:decimal(11,8);not null" json:"long"`

	GeoHash string `gorm:"type:varchar(12);index:idx_stores_geohash" json:"geoHash"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

func NewStore(name string, address string, lat float64, long float64) Store {
	return Store{
		Name:          name,
		StreetAddress: address,
		Lat:           lat,
		Long:          long,
	}
}

// ---------------------------------------------------------------------------

type Category struct {
	ID   uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
	Slug string `gorm:"type:varchar(100);uniqueIndex;not null" json:"slug"`
}

// ---------------------------------------------------------------------------

type Product struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`

	Name      string  `gorm:"type:varchar(255);not null" json:"name"`
	Price     float64 `gorm:"type:decimal(12,2);not null" json:"price"`
	PriceUnit string  `gorm:"type:varchar(20);not null" json:"priceUnit"`

	Quantity     *float64 `gorm:"type:decimal(12,2)" json:"quantity"`
	QuantityUnit *string  `gorm:"type:varchar(20)" json:"quantityUnit"`
	Desc         string   `gorm:"type:text" json:"desc"`

	BusinessID uint      `gorm:"index;not null" json:"businessId"`
	Business   *Business `gorm:"foreignKey:BusinessID;constraint:OnDelete:CASCADE;" json:"business,omitempty"`

	SearchTokens []string   `gorm:"type:jsonb;serializer:json" json:"-"`
	Categories   []Category `gorm:"many2many:product_categories;" json:"categories,omitempty"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

func NewProduct(name string, price float64, priceUnit string, desc string, quantity *float64, quantityUnit *string, searchTokens []string) Product {
	return Product{
		Name:         name,
		Price:        price,
		PriceUnit:    priceUnit,
		Desc:         desc,
		Quantity:     quantity,
		QuantityUnit: quantityUnit,
		SearchTokens: searchTokens,
	}
}

// ---------------------------------------------------------------------------

type InventoryProduct struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`

	StoreID uint   `gorm:"uniqueIndex:idx_store_product;not null" json:"storeId"`
	Store   *Store `gorm:"foreignKey:StoreID;constraint:OnDelete:CASCADE;" json:"store,omitempty"`

	ProductID uint     `gorm:"uniqueIndex:idx_store_product;not null" json:"productId"`
	Product   *Product `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;" json:"product,omitempty"`

	Count *int `gorm:"default:null" json:"count"`

	Available *bool `gorm:"default:true" json:"available"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

func NewInvProduct(storeID uint, productID uint, count *int, available *bool) InventoryProduct {
	return InventoryProduct{
		StoreID:   storeID,
		ProductID: productID,
		Count:     count,
		Available: available,
	}
}

// ---------------------------------------------------------------------------

type BusinessAccess struct {
	ID     uint `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID uint `gorm:"index;unique;not null" json:"userId"`

	BusinessID uint      `gorm:"index;not null" json:"businessId"`
	Business   *Business `gorm:"foreignKey:BusinessID;constraint:OnDelete:CASCADE;" json:"business,omitempty"`

	Role jwtutil.UserRole `gorm:"type:varchar(50);default:'owner';not null" json:"role"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewBusinessAccess(userID uint, businessID uint, role jwtutil.UserRole) BusinessAccess {
	return BusinessAccess{
		UserID:     userID,
		BusinessID: businessID,
		Role:       role,
	}
}

// ---------------------------------------------------------------------------

type ProductToken struct {
	Token      string    `gorm:"primaryKey;type:varchar(100)" json:"token"`
	CategoryID uint      `gorm:"primaryKey" json:"categoryId"`
	Category   *Category `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE" json:"category,omitempty"`
	Count      int       `gorm:"not null;default:1" json:"count"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// ---------------------------------------------------------------------------

// BusinessAccount stores credentials for business owners who authenticate
// directly with the vendor service, independent of the user service.
type BusinessAccount struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Email        string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"`
	BusinessID   uint      `gorm:"uniqueIndex;not null" json:"businessId"`
	Business     *Business `gorm:"foreignKey:BusinessID;constraint:OnDelete:CASCADE;" json:"business,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// ---------------------------------------------------------------------------

// SearchSyncOutbox is the transactional outbox for propagating inventory changes
// to the Search Service. Written in the same DB transaction as the triggering write.
type SearchSyncOutbox struct {
	ID          int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	EventType   string     `gorm:"type:varchar(20);not null" json:"eventType"` // "upsert" | "delete"
	Payload     []byte     `gorm:"type:jsonb;not null" json:"payload"`
	Processed   bool       `gorm:"default:false;index" json:"processed"`
	ProcessedAt *time.Time `json:"processedAt,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
}

// OutboxPayload is the JSON stored in SearchSyncOutbox.Payload.
type OutboxPayload struct {
	EventType     string   `json:"event_type"`
	InvProductID  uint     `json:"inv_product_id"`
	ProductID     uint     `json:"product_id,omitempty"`
	BusinessID    uint     `json:"business_id,omitempty"`
	ProductName   string   `json:"product_name,omitempty"`
	Price         float64  `json:"price,omitempty"`
	PriceUnit     string   `json:"price_unit,omitempty"`
	Quantity      *float64 `json:"quantity,omitempty"`
	QuantityUnit  *string  `json:"quantity_unit,omitempty"`
	Desc          string   `json:"desc,omitempty"`
	SearchTokens  []string `json:"search_tokens,omitempty"`
	Categories    []CatRef `json:"categories,omitempty"`
	StoreID       uint     `json:"store_id,omitempty"`
	StoreName     string   `json:"store_name,omitempty"`
	StreetAddress string   `json:"street_address,omitempty"`
	Lat           float64  `json:"lat,omitempty"`
	Long          float64  `json:"long,omitempty"`
	GeoHash       string   `json:"geo_hash,omitempty"`
	Available     bool     `json:"available"`
}

type CatRef struct {
	ID uint `json:"id"`
}
