package models

import (
	"time"

	"github.com/atharvyadav96k/spotnearr/pkg/jwtutil"
)

type BusinessAccess struct {
	ID     uint  `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID uint  `gorm:"index;unique;not null" json:"userId"`

	BusinessID uint      `gorm:"index;not null" json:"businessId"`
	Business   *Business `gorm:"foreignKey:BusinessID;constraint:OnDelete:CASCADE;" json:"business,omitempty"`

	Role jwtutil.UserRole `gorm:"type:varchar(50);default:'owner';not null" json:"role"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewBusinessAccess(userId uint, businessId uint, role jwtutil.UserRole) BusinessAccess {
	return BusinessAccess{
		UserID:     userId,
		BusinessID: businessId,
		Role:       role,
	}
}
