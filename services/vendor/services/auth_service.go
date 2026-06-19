package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/atharvyadav96k/spotnearr/pkg/httputil"
	"github.com/atharvyadav96k/spotnearr/pkg/jwtutil"
	"github.com/atharvyadav96k/spotnearr/pkg/validate"
	vendormodel "github.com/Developer-Aadesh/spotnearr-database/vendordb"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ─── DTOs ────────────────────────────────────────────────────────────────────

type RegisterBusinessDTO struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Desc     string `json:"desc"`
}

func (d *RegisterBusinessDTO) Validate() error {
	if strings.TrimSpace(d.Name) == "" {
		return fmt.Errorf("name is required")
	}
	if !validate.Email(d.Email) {
		return fmt.Errorf("invalid email address")
	}
	if !validate.Phone(d.Phone) {
		return fmt.Errorf("invalid phone number")
	}
	if ok, msg := validate.Password(d.Password); !ok {
		return fmt.Errorf("%s", msg)
	}
	return nil
}

type LoginBusinessDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (d *LoginBusinessDTO) Validate() error {
	if strings.TrimSpace(d.Email) == "" {
		return fmt.Errorf("email is required")
	}
	if strings.TrimSpace(d.Password) == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}

type RefreshDTO struct {
	RefreshToken string `json:"refresh_token"`
}

func (d *RefreshDTO) Validate() error {
	if strings.TrimSpace(d.RefreshToken) == "" {
		return fmt.Errorf("refresh_token is required")
	}
	return nil
}

// ─── Service ─────────────────────────────────────────────────────────────────

type AuthService struct {
	db        *gorm.DB
	jwtSecret string
}

func NewAuthService(db *gorm.DB, jwtSecret string) *AuthService {
	return &AuthService{db: db, jwtSecret: jwtSecret}
}

func (a *AuthService) Register(dto *RegisterBusinessDTO) httputil.Res {
	hash, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return httputil.NewResponse("Failed to process credentials", http.StatusInternalServerError, nil)
	}

	var account vendormodel.BusinessAccount
	err = a.db.WithContext(context.Background()).Transaction(func(tx *gorm.DB) error {
		biz := vendormodel.NewBusiness(dto.Name, dto.Email, dto.Phone, dto.Desc)
		if err := tx.Create(biz).Error; err != nil {
			return err
		}
		account = vendormodel.BusinessAccount{
			Email:        dto.Email,
			PasswordHash: string(hash),
			BusinessID:   biz.ID,
		}
		return tx.Create(&account).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return httputil.NewResponse("Email or phone already in use", http.StatusConflict, nil)
		}
		return httputil.NewResponse(err.Error(), http.StatusBadRequest, nil)
	}

	return a.issueTokens(account)
}

func (a *AuthService) Login(dto *LoginBusinessDTO) httputil.Res {
	var account vendormodel.BusinessAccount
	if err := a.db.WithContext(context.Background()).
		Where("email = ?", dto.Email).
		First(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return httputil.NewResponse("Invalid credentials", http.StatusUnauthorized, nil)
		}
		return httputil.NewResponse("Login failed", http.StatusInternalServerError, nil)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(dto.Password)); err != nil {
		return httputil.NewResponse("Invalid credentials", http.StatusUnauthorized, nil)
	}

	return a.issueTokens(account)
}

func (a *AuthService) Refresh(dto *RefreshDTO) httputil.Res {
	claims, err := jwtutil.ValidateToken(dto.RefreshToken, a.jwtSecret, jwtutil.TypeRefreshToken)
	if err != nil {
		return httputil.NewResponse("Invalid or expired refresh token", http.StatusUnauthorized, nil)
	}

	accessToken, expiry, err := jwtutil.GenerateAccessToken(claims.UserId, claims.BusinessId, a.jwtSecret, jwtutil.RoleAdmin)
	if err != nil {
		return httputil.NewResponse("Failed to generate token", http.StatusInternalServerError, nil)
	}

	return httputil.NewResponse("Token refreshed", http.StatusOK, jwtutil.TokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: expiry.Unix(),
	})
}

func (a *AuthService) issueTokens(account vendormodel.BusinessAccount) httputil.Res {
	bizID := account.BusinessID
	accessToken, accessExpiry, err := jwtutil.GenerateAccessToken(account.ID, &bizID, a.jwtSecret, jwtutil.RoleAdmin)
	if err != nil {
		return httputil.NewResponse("Failed to generate access token", http.StatusInternalServerError, nil)
	}
	refreshToken, refreshExpiry, err := jwtutil.GenerateRefreshToken(account.ID, &bizID, a.jwtSecret, jwtutil.RoleAdmin)
	if err != nil {
		return httputil.NewResponse("Failed to generate refresh token", http.StatusInternalServerError, nil)
	}
	return httputil.NewResponse("Success", http.StatusOK, jwtutil.TokenResponse{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessExpiry.Unix(),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshExpiry.Unix(),
	})
}
