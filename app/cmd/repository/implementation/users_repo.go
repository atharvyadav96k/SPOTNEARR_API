package implementation

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/atharvyadav96k/SPOTNEARR_API/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) Register(ctx context.Context, user *models.User) error {
	var existingUser models.User
	err := u.db.WithContext(ctx).
		Where("email = ? OR phone = ?", user.Email, user.Phone).
		First(&existingUser).Error

	if err == nil {
		if user.Email != nil && existingUser.Email != nil && *user.Email == *existingUser.Email {
			return fmt.Errorf("User with this email already exists")
		}
		if user.Phone != nil && existingUser.Phone != nil && *user.Phone == *existingUser.Phone {
			return fmt.Errorf("User with this phone already exists")
		}
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if err := u.db.WithContext(ctx).Create(user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			errStr := err.Error()
			if strings.Contains(errStr, "email") {
				return fmt.Errorf("User with this email already exists")
			}
			if strings.Contains(errStr, "phone") {
				return fmt.Errorf("User with this phone already exists")
			}
		}
		return err
	}
	return nil
}

func (u *UserRepository) SetRefreshToken(ctx context.Context, userID uint, token string) error {
	db := u.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Update("refresh_token", token)

	if db.Error != nil {
		return db.Error
	}
	if db.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (u *UserRepository) RemoveRefreshToken(ctx context.Context, userID uint) error {
	db := u.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Update("refresh_token", "")

	if db.Error != nil {
		return db.Error
	}
	if db.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (u *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := u.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("User not found with email %v", email)
		}
		return nil, err
	}

	if !user.IsActive {
		return nil, fmt.Errorf("Your account is deactivated. Please request system permission to reactivate")
	}

	return &user, nil
}

func (u *UserRepository) GetById(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	db := u.db.WithContext(ctx).Where("id = ?", id).First(&user)
	if db.Error != nil {
		return nil, db.Error
	}
	if db.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &user, nil
}

func (u *UserRepository) GetByPhone(ctx context.Context, phone string) (*models.User, error) {
	var user models.User
	if err := u.db.WithContext(ctx).Where("phone = ?", phone).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("User not found with phone %v", phone)
		}
		return nil, err
	}

	if !user.IsActive {
		return nil, fmt.Errorf("Your account is deactivated. Please request system permission to reactivate")
	}

	return &user, nil
}

func (u *UserRepository) UpdatePasswordWithEmail(ctx context.Context, email string, hashedPassword string) error {
	db := u.db.WithContext(ctx).Model(&models.User{}).
		Where("email = ? AND is_active = ?", email, true).
		Update("password_hash", hashedPassword)

	if db.Error != nil {
		return fmt.Errorf("Failed to update the password")
	}

	if db.RowsAffected == 0 {
		return fmt.Errorf("Cannot update password: user does not exist or is inactive")
	}

	return nil
}

func (u *UserRepository) UpdatePasswordWithUserId(ctx context.Context, userId uint, hasedPassword string) error {
	db := u.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ? AND is_active = ?", userId, true).
		Update("password_hash", hasedPassword)
	if db.Error != nil {
		return fmt.Errorf("Failed to update the password")
	}
	if db.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (u *UserRepository) FreezeAccount(ctx context.Context, email string) error {
	db := u.db.WithContext(ctx).Model(&models.User{}).
		Where("email = ?", email).
		Update("is_active", false)

	if db.Error != nil {
		return db.Error
	}

	if db.RowsAffected == 0 {
		return fmt.Errorf("Account not found with this email %v", email)
	}

	return nil
}

func (u *UserRepository) UnfreezeAccount(ctx context.Context, email string) error {
	db := u.db.WithContext(ctx).Model(&models.User{}).
		Where("email = ?", email).
		Updates(map[string]interface{}{
			"is_active": true,
		})

	if db.Error != nil {
		return fmt.Errorf("failed to reactivate account: %w", db.Error)
	}
	if db.RowsAffected == 0 {
		return fmt.Errorf("account not found with this email %v", email)
	}
	return nil
}

// func (u *UserRepository) SetBusinessId(ctx context.Context, email string, businessId uint) error {
// 	db := u.db.WithContext(ctx).Model(&models.User).
// 		Where("email = ?", email).
// 		Update("business_id", businessId)

// }
