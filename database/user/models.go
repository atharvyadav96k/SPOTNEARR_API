package user

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleOwner UserRole = "owner"
)

type User struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	FullName string `gorm:"type:varchar(255);not null" json:"fullName"`

	Email *string `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Phone *string `gorm:"type:varchar(20);uniqueIndex;not null" json:"phone"`

	PasswordHash    *string `gorm:"type:text;not null" json:"-"`
	IsVerifiedEmail bool    `gorm:"default:false;not null" json:"isVerifiedEmail"`
	IsVerifiedPhone bool    `gorm:"default:false;not null" json:"isVerifiedPhone"`
	IsActive        bool    `gorm:"default:true;not null" json:"isActive"`

	PushToken    *string `gorm:"type:text" json:"-"`
	RefreshToken *string `gorm:"type:text" json:"-"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

func NewUser(name string, email string, phone string, password string) *User {
	return &User{
		FullName:     name,
		Email:        &email,
		Phone:        &phone,
		PasswordHash: &password,
	}
}

// ---------------------------------------------------------------------------

type ClaimStatus string

const (
	ClaimStatusPending  ClaimStatus = "pending"
	ClaimStatusAccepted ClaimStatus = "accepted"
	ClaimStatusRejected ClaimStatus = "rejected"
)

type Claim struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`

	UserID uint  `gorm:"index;not null" json:"userId"`
	User   *User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"user,omitempty"`

	InventoryProductID uint `gorm:"index;not null" json:"inventoryProductId"`

	Status ClaimStatus `gorm:"type:varchar(20);default:'pending';not null" json:"status"`

	// Frozen snapshot of the inventory product at claim time so GET /claims
	// never requires a cross-service call.
	InvProductSnapshot json.RawMessage `gorm:"type:jsonb" json:"product,omitempty"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

func NewClaim(userID uint, invProductID uint, snapshot json.RawMessage) Claim {
	return Claim{
		UserID:             userID,
		InventoryProductID: invProductID,
		Status:             ClaimStatusPending,
		InvProductSnapshot: snapshot,
	}
}

// ---------------------------------------------------------------------------

type ReviewTarget string

const (
	ReviewTargetBusiness  ReviewTarget = "business"
	ReviewTargetProduct   ReviewTarget = "product"
	ReviewTargetOffer     ReviewTarget = "offer"
	ReviewTargetSpotlight ReviewTarget = "spotlight"
)

type Review struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`

	UserID uint  `gorm:"index;not null" json:"userId"`
	User   *User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"user,omitempty"`

	TargetType ReviewTarget `gorm:"type:varchar(20);not null;index:idx_review_target" json:"targetType"`
	TargetID   uint         `gorm:"not null;index:idx_review_target" json:"targetId"`

	Stars   uint8  `gorm:"not null" json:"stars"`
	Comment string `gorm:"type:text" json:"comment"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

// ---------------------------------------------------------------------------

// BusinessFollow records that a user follows a business.
// Composite PK prevents duplicates. No soft-delete — a delete is a real unfollow.
type BusinessFollow struct {
	UserID     uint      `gorm:"primaryKey;not null" json:"userId"`
	BusinessID uint      `gorm:"primaryKey;not null" json:"businessId"`
	CreatedAt  time.Time `json:"createdAt"`
}

// ---------------------------------------------------------------------------

// ProductLike records that a user liked an inventory product.
type ProductLike struct {
	UserID       uint      `gorm:"primaryKey;not null" json:"userId"`
	InvProductID uint      `gorm:"primaryKey;not null" json:"invProductId"`
	CreatedAt    time.Time `json:"createdAt"`
}

// ProductSave records that a user saved an inventory product.
type ProductSave struct {
	UserID       uint      `gorm:"primaryKey;not null" json:"userId"`
	InvProductID uint      `gorm:"primaryKey;not null" json:"invProductId"`
	CreatedAt    time.Time `json:"createdAt"`
}

// ---------------------------------------------------------------------------

// SpotlightLike records that a user liked a spotlight post.
type SpotlightLike struct {
	UserID      uint      `gorm:"primaryKey;not null" json:"userId"`
	SpotlightID uint      `gorm:"primaryKey;not null" json:"spotlightId"`
	CreatedAt   time.Time `json:"createdAt"`
}

// SpotlightSave records that a user saved a spotlight post.
type SpotlightSave struct {
	UserID      uint      `gorm:"primaryKey;not null" json:"userId"`
	SpotlightID uint      `gorm:"primaryKey;not null" json:"spotlightId"`
	CreatedAt   time.Time `json:"createdAt"`
}
