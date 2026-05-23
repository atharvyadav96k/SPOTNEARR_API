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
	ID uint `gorm:"primaryKey;autoIncrement"`

	FullName string `gorm:"type:varchar(255);not null"`

	Email *string `gorm:"type:varchar(255);unique;index"`
	Phone *string `gorm:"type:varchar(20);unique;index"`

	PasswordHash *string `gorm:"type:text;not null"`

	Role UserRole `gorm:"type:varchar(50);default:'user';not null"`

	IsVerifiedEmail bool `gorm:"default:false;not null"`
	IsVerifiedPhone bool `gorm:"default:false;not null"`
	IsActive        bool `gorm:"default:true;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
