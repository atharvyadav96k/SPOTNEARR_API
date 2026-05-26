package models

import (
	"time"
)

const (
	RoleUser     UserRole = "user"
	RoleAdmin    UserRole = "admin"
	RoleSubAdmin UserRole = "sub-admin"
)

type BusinessAccess struct {
	ID     uint  `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID uint  `gorm:"index;unique;not null" json:"userId"`
	User   *User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"user,omitempty"`

	BusinessID uint      `gorm:"index;not null" json:"businessId"`
	Business   *Business `gorm:"foreignKey:BusinessID;constraint:OnDelete:CASCADE;" json:"business,omitempty"`

	Role UserRole `gorm:"type:varchar(50);default:'owner';not null" json:"role"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewBusinessAccess(userId uint, businessId uint, role UserRole) BusinessAccess {
	return BusinessAccess{
		UserID:     userId,
		BusinessID: businessId,
		Role:       role,
	}
}
