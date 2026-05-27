package repository

import (
	"context"

	"github.com/atharvyadav96k/SPOTNEARR_API/models"
)

type IUserRepository interface {
	Register(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetById(ctx context.Context, id uint) (*models.User, error)
	GetByPhone(ctx context.Context, phone string) (*models.User, error)
	SetRefreshToken(ctx context.Context, userID uint, token string) error
	RemoveRefreshToken(ctx context.Context, userID uint) error
	UpdatePasswordWithEmail(ctx context.Context, email string, hashedPassword string) error
	UpdatePasswordWithUserId(ctx context.Context, userId uint, hashedPassword string) error
	FreezeAccount(ctx context.Context, email string) error
	UnfreezeAccount(ctx context.Context, email string) error
}
