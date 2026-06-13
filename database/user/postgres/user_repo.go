package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Developer-Aadesh/spotnearr-database/user"
	"gorm.io/gorm"
)

type UserRepository struct{ db *gorm.DB }

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) Register(ctx context.Context, usr *user.User) error {
	var existing user.User
	err := u.db.WithContext(ctx).
		Where("email = ? OR phone = ?", usr.Email, usr.Phone).
		First(&existing).Error

	if err == nil {
		if usr.Email != nil && existing.Email != nil && *usr.Email == *existing.Email {
			return fmt.Errorf("user with this email already exists")
		}
		if usr.Phone != nil && existing.Phone != nil && *usr.Phone == *existing.Phone {
			return fmt.Errorf("user with this phone already exists")
		}
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if err := u.db.WithContext(ctx).Create(usr).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			s := err.Error()
			if strings.Contains(s, "email") {
				return fmt.Errorf("user with this email already exists")
			}
			if strings.Contains(s, "phone") {
				return fmt.Errorf("user with this phone already exists")
			}
		}
		return err
	}
	return nil
}

func (u *UserRepository) SetRefreshToken(ctx context.Context, userID uint, token string) error {
	db := u.db.WithContext(ctx).Model(&user.User{}).Where("id = ?", userID).Update("refresh_token", token)
	if db.Error != nil {
		return db.Error
	}
	if db.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (u *UserRepository) RemoveRefreshToken(ctx context.Context, userID uint) error {
	db := u.db.WithContext(ctx).Model(&user.User{}).Where("id = ?", userID).Update("refresh_token", "")
	if db.Error != nil {
		return db.Error
	}
	if db.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (u *UserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	var usr user.User
	if err := u.db.WithContext(ctx).Where("email = ?", email).First(&usr).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found with email %v", email)
		}
		return nil, err
	}
	if !usr.IsActive {
		return nil, fmt.Errorf("your account is deactivated, please request system permission to reactivate")
	}
	return &usr, nil
}

func (u *UserRepository) GetById(ctx context.Context, id uint) (*user.User, error) {
	var usr user.User
	db := u.db.WithContext(ctx).Where("id = ?", id).First(&usr)
	if db.Error != nil {
		return nil, db.Error
	}
	if db.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &usr, nil
}

func (u *UserRepository) GetByPhone(ctx context.Context, phone string) (*user.User, error) {
	var usr user.User
	if err := u.db.WithContext(ctx).Where("phone = ?", phone).First(&usr).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found with phone %v", phone)
		}
		return nil, err
	}
	if !usr.IsActive {
		return nil, fmt.Errorf("your account is deactivated, please request system permission to reactivate")
	}
	return &usr, nil
}

func (u *UserRepository) UpdatePasswordWithEmail(ctx context.Context, email string, hashedPassword string) error {
	db := u.db.WithContext(ctx).Model(&user.User{}).
		Where("email = ? AND is_active = ?", email, true).
		Update("password_hash", hashedPassword)
	if db.Error != nil {
		return fmt.Errorf("failed to update the password")
	}
	if db.RowsAffected == 0 {
		return fmt.Errorf("cannot update password: user does not exist or is inactive")
	}
	return nil
}

func (u *UserRepository) UpdatePasswordWithUserId(ctx context.Context, userId uint, hashedPassword string) error {
	db := u.db.WithContext(ctx).Model(&user.User{}).
		Where("id = ? AND is_active = ?", userId, true).
		Update("password_hash", hashedPassword)
	if db.Error != nil {
		return fmt.Errorf("failed to update the password")
	}
	if db.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (u *UserRepository) UpdateProfile(ctx context.Context, userID uint, fullName string) error {
	db := u.db.WithContext(ctx).Model(&user.User{}).Where("id = ? AND is_active = ?", userID, true).Update("full_name", fullName)
	if db.Error != nil {
		return db.Error
	}
	if db.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (u *UserRepository) FreezeAccountByID(ctx context.Context, userID uint) error {
	db := u.db.WithContext(ctx).Model(&user.User{}).Where("id = ?", userID).Update("is_active", false)
	if db.Error != nil {
		return db.Error
	}
	if db.RowsAffected == 0 {
		return fmt.Errorf("account not found with id %d", userID)
	}
	return nil
}

func (u *UserRepository) FreezeAccount(ctx context.Context, email string) error {
	db := u.db.WithContext(ctx).Model(&user.User{}).Where("email = ?", email).Update("is_active", false)
	if db.Error != nil {
		return db.Error
	}
	if db.RowsAffected == 0 {
		return fmt.Errorf("account not found with email %v", email)
	}
	return nil
}

func (u *UserRepository) UnfreezeAccount(ctx context.Context, email string) error {
	db := u.db.WithContext(ctx).Model(&user.User{}).
		Where("email = ?", email).
		Updates(map[string]interface{}{"is_active": true})
	if db.Error != nil {
		return fmt.Errorf("failed to reactivate account: %w", db.Error)
	}
	if db.RowsAffected == 0 {
		return fmt.Errorf("account not found with email %v", email)
	}
	return nil
}
