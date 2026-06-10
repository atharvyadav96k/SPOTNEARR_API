package user

import "context"

type IUserRepository interface {
	Register(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetById(ctx context.Context, id uint) (*User, error)
	GetByPhone(ctx context.Context, phone string) (*User, error)
	SetRefreshToken(ctx context.Context, userID uint, token string) error
	RemoveRefreshToken(ctx context.Context, userID uint) error
	UpdatePasswordWithEmail(ctx context.Context, email string, hashedPassword string) error
	UpdatePasswordWithUserId(ctx context.Context, userId uint, hashedPassword string) error
	FreezeAccount(ctx context.Context, email string) error
	UnfreezeAccount(ctx context.Context, email string) error
}

type IClaimRepository interface {
	Create(ctx context.Context, claim Claim) (Claim, error)
	GetByID(ctx context.Context, claimID uint) (Claim, error)
	GetByUserID(ctx context.Context, userID uint) ([]Claim, error)
	GetByUserAndInvProduct(ctx context.Context, userID uint, invProductID uint) (Claim, error)
	Delete(ctx context.Context, claimID uint, userID uint) error
}

type IReviewRepository interface {
	Add(ctx context.Context, review Review) (Review, error)
	Update(ctx context.Context, userID uint, targetType ReviewTarget, targetID uint, stars uint8, comment string) (Review, error)
	Delete(ctx context.Context, userID uint, targetType ReviewTarget, targetID uint) error
	GetByTarget(ctx context.Context, targetType ReviewTarget, targetID uint) ([]Review, error)
	GetByUserAndTarget(ctx context.Context, userID uint, targetType ReviewTarget, targetID uint) (Review, error)
}
