package models

import (
	"time"

	"gorm.io/gorm"
)

type UserRole string

type User struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	FullName string `gorm:"type:varchar(255);not null" json:"fullName"`

	Email *string `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Phone *string `gorm:"type:varchar(20);uniqueIndex;not null" json:"phone"`

	PasswordHash    *string `gorm:"type:text;not null" json:"-"`
	IsVerifiedEmail bool    `gorm:"default:false;not null" json:"isVerifiedEmail"`
	IsVerifiedPhone bool    `gorm:"default:false;not null" json:"isVerifiedPhone"`
	IsActive        bool    `gorm:"default:true;not null" json:"isActive"`

	PushToken *string `gorm:"type:text" json:"-"`

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
