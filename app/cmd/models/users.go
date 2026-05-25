package models

import (
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
)

type User struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	FullName string `gorm:"type:varchar(255);not null" json:"fullName"`

	Email *string `gorm:"type:varchar(255);uniqueIndex" json:"email"`
	Phone *string `gorm:"type:varchar(20);uniqueIndex" json:"phone"`

	PasswordHash *string  `gorm:"type:text;not null" json:"-"`
	Role         UserRole `gorm:"type:varchar(50);default:'user';not null" json:"role"`

	BusinessID *uint     `gorm:"index" json:"businessId"`
	Business   *Business `gorm:"foreignKey:BusinessID;constraint:OnDelete:SET NULL;" json:"business,omitempty"`

	RefreshToken string `gorm:"type:varchar(255);" json:"-"`

	IsVerifiedEmail bool `gorm:"default:false;not null" json:"isVerifiedEmail"`
	IsVerifiedPhone bool `gorm:"default:false;not null" json:"isVerifiedPhone"`
	IsActive        bool `gorm:"default:true;not null" json:"isActive"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}
