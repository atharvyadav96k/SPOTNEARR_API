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
	UpdateProfile(ctx context.Context, userID uint, fullName string) error
	FreezeAccount(ctx context.Context, email string) error
	FreezeAccountByID(ctx context.Context, userID uint) error
	UnfreezeAccount(ctx context.Context, email string) error
}

type IClaimRepository interface {
	Create(ctx context.Context, claim Claim) (Claim, error)
	GetByID(ctx context.Context, claimID uint) (Claim, error)
	GetByUserID(ctx context.Context, userID uint) ([]Claim, error)
	GetByUserAndInvProduct(ctx context.Context, userID uint, invProductID uint) (Claim, error)
	GetByProductIDs(ctx context.Context, invProductIDs []uint) ([]Claim, error)
	UpdateStatus(ctx context.Context, claimID uint, status ClaimStatus) error
	Delete(ctx context.Context, claimID uint, userID uint) error
}

type IReviewRepository interface {
	Add(ctx context.Context, review Review) (Review, error)
	Update(ctx context.Context, userID uint, targetType ReviewTarget, targetID uint, stars uint8, comment string) (Review, error)
	Delete(ctx context.Context, userID uint, targetType ReviewTarget, targetID uint) error
	GetByTarget(ctx context.Context, targetType ReviewTarget, targetID uint) ([]Review, error)
	GetByUserAndTarget(ctx context.Context, userID uint, targetType ReviewTarget, targetID uint) (Review, error)
}

type IBusinessFollowRepository interface {
	Follow(ctx context.Context, userID, businessID uint) error
	Unfollow(ctx context.Context, userID, businessID uint) error
	IsFollowing(ctx context.Context, userID, businessID uint) (bool, error)
	GetFollowedBusinessIDs(ctx context.Context, userID uint) ([]uint, error)
}

type IProductEngagementRepository interface {
	Like(ctx context.Context, userID, invProductID uint) error
	Unlike(ctx context.Context, userID, invProductID uint) error
	Save(ctx context.Context, userID, invProductID uint) error
	Unsave(ctx context.Context, userID, invProductID uint) error
	GetLikedByUser(ctx context.Context, userID uint) ([]uint, error)
	GetSavedByUser(ctx context.Context, userID uint) ([]uint, error)
}

type ISpotlightEngagementRepository interface {
	Like(ctx context.Context, userID, spotlightID uint) error
	Unlike(ctx context.Context, userID, spotlightID uint) error
	Save(ctx context.Context, userID, spotlightID uint) error
	Unsave(ctx context.Context, userID, spotlightID uint) error
	GetLikedByUser(ctx context.Context, userID uint) ([]uint, error)
	GetSavedByUser(ctx context.Context, userID uint) ([]uint, error)
}
