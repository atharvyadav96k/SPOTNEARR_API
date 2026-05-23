package repository

import (
	"context"

	"github.com/atharvyadav96k/SPOTNEARR_API/models"
)

type IUserRepository interface {
	Register(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByPhone(ctx context.Context, phone string) (*models.User, error)
	UpdatePasswordWithEmail(ctx context.Context, email string, hashedPassword string) error
	FreezeAccount(ctx context.Context, email string) error
	UnfreezeAccount(ctx context.Context, email string) error
}
