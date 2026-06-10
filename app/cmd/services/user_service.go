package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/atharvyadav96k/SPOTNEARR_API/auth"
	"github.com/atharvyadav96k/SPOTNEARR_API/config"
	"github.com/atharvyadav96k/SPOTNEARR_API/connections/cache"
	"github.com/atharvyadav96k/SPOTNEARR_API/models"
	"github.com/atharvyadav96k/SPOTNEARR_API/utils"
	"github.com/atharvyadav96k/SPOTNEARR_API/utils/response"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	base_service
	httpClient *http.Client
}

func NewUserService(db *gorm.DB, cache *cache.Cache) *UserService {
	return &UserService{
		base_service: NewBaseService(db, cache),
		httpClient:   &http.Client{Timeout: 500 * time.Millisecond},
	}
}

// fetchBusinessAccess calls the Vendor Service to get the user's business_id + role.
// Returns nil, RoleUser if the user has no business or if the call fails.
func (u *UserService) fetchBusinessAccess(userID uint) (*uint, models.UserRole) {
	url := fmt.Sprintf("%s/internal/users/%d/access", config.C.VendorServiceURL, userID)
	resp, err := u.httpClient.Get(url)
	if err != nil {
		log.Printf("user: vendor access lookup failed for user %d: %v", userID, err)
		return nil, models.RoleUser
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, models.RoleUser // user has no business yet
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("user: vendor access returned %d for user %d", resp.StatusCode, userID)
		return nil, models.RoleUser
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, models.RoleUser
	}
	var result struct {
		BusinessID uint             `json:"business_id"`
		Role       models.UserRole  `json:"role"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, models.RoleUser
	}
	return &result.BusinessID, result.Role
}

func (u *UserService) RegisterUser(user *models.User) response.Res {
	if user == nil {
		return u.ResponseBadRequest("Failed to get user")
	}
	hashedPassword, err := auth.GenerateHashedPassword(*user.PasswordHash)
	if err != nil {
		return u.ResponseBadRequest("Failed to process password")
	}
	user.PasswordHash = &hashedPassword
	if err := u.RepoUser().Register(context.Background(), user); err != nil {
		return u.ResponseBadRequest(err.Error())
	}
	return u.ResponseCreated("User registered successfully", nil)
}

func (u *UserService) Login(email string, password string) response.Res {
	ctx := context.Background()
	user, err := u.RepoUser().GetByEmail(ctx, email)
	if err != nil {
		return u.ResponseBadRequest("Invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(password))
	if err != nil {
		log.Default().Println(err)
		return u.ResponseBadRequest("Invalid email or password")
	}

	currentBusinessID, userRole := u.fetchBusinessAccess(user.ID)

	var refreshToken string
	var refreshExpiry time.Time

	cachedToken, err := u.Cache().GetRefreshTokenSession().GetRefreshTokenSession(user.ID)
	if err == nil {
		if rtClaims, err := auth.ValidateToken(cachedToken, config.C.JWTSecret, auth.TypeRefreshToken); err == nil {
			refreshToken = cachedToken
			refreshExpiry = rtClaims.ExpiresAt.Time
		}
	}

	if refreshToken == "" {
		refreshToken, refreshExpiry, err = auth.GenerateRefreshToken(user.ID, currentBusinessID, config.C.JWTSecret, userRole)
		if err != nil {
			log.Default().Println("Failed to generate refresh token:", err)
			return u.ResponseInternalServer("Failed to login")
		}

		err = u.Cache().GetRefreshTokenSession().NewRefreshToken(user.ID, refreshToken, 30*24*time.Hour)
		if err != nil {
			log.Default().Println("Failed to save refresh token:", err)
			return u.ResponseInternalServer("Failed to login")
		}
	}

	accessToken, accessExpiry, err := auth.GenerateAccessToken(user.ID, currentBusinessID, config.C.JWTSecret, userRole)
	if err != nil {
		log.Default().Println("Failed to generate access token:", err)
		return u.ResponseInternalServer("Failed to login")
	}

	return u.ResponseOK("Logged in successfully", auth.TokenResponse{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessExpiry.Unix(),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshExpiry.Unix(),
	})
}

func (u *UserService) Refresh(claims auth.UserClaims, refreshToken string) response.Res {
	cachedToken, err := u.Cache().GetRefreshTokenSession().GetRefreshTokenSession(claims.UserId)
	if err != nil || cachedToken != refreshToken {
		log.Default().Println("Refresh token not found in cache or mismatch")
		return u.ResponseUnauthorized()
	}

	accessToken, accessExpiry, err := auth.GenerateAccessToken(claims.UserId, claims.BusinessId, config.C.JWTSecret, claims.UserRole)
	if err != nil {
		log.Default().Println("Failed to generate access token")
		return u.ResponseUnauthorized()
	}

	return u.ResponseOK("New access token", auth.TokenResponse{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessExpiry.Unix(),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: claims.ExpiresAt.Time.Unix(),
	})
}

func (u *UserService) DismissRefreshToken(id uint) response.Res {
	if err := u.Cache().GetRefreshTokenSession().InvalidateRefreshToken(id); err != nil {
		log.Default().Println("Failed to invalidate refresh token:", err)
		return u.ResponseInternalServer("Failed to logout")
	}
	return u.ResponseOK("Successfully logged out from all devices", nil)
}

// InvalidateRefreshToken drops the stored refresh token for a user.
// Used by the internal endpoint called by the Vendor Service after business registration.
func (u *UserService) InvalidateRefreshToken(id uint) error {
	return u.Cache().GetRefreshTokenSession().InvalidateRefreshToken(id)
}

func (u *UserService) AddUserToBusiness(userID uint, businessId uint, ownerId uint) response.Res {
	// if err := u.RepoUser().
	return response.Res{}
}

func (u *UserService) GetUserProfile(userId uint) response.Res {
	user, err := u.RepoUser().GetById(context.Background(), userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return u.ResponseNotFound("User not found")
		}
		return u.ResponseBadRequest("Failed to get user profile")
	}
	return u.ResponseOK("User profile", user)
}

func (u *UserService) SessionNotification(email string) response.Res {
	expireTime := time.Hour * 1
	session, err := utils.GenerateSession(email, expireTime)
	if err != nil {
		return u.ResponseInternalServer("Failed to send the email.")
	}
	_, err = u.Cache().GetPasswordSessions().NewPasswordSession(context.Background(), email, expireTime, session)
	if err != nil {
		return u.ResponseBadRequest("Failed to send the email.")
	}
	sessionLink := fmt.Sprintf("session=%s", session)
	log.Default().Println(sessionLink)
	return u.ResponseOK("", sessionLink)
}

func (u *UserService) UpdatePassword(password string, session string) response.Res {
	ctx := context.Background()

	email, err := u.Cache().GetPasswordSessions().GetPasswordSession(ctx, session)
	if err != nil {
		return u.ResponseInternalServer("Failed to update password")
	}
	if email == "" {
		return u.ResponseBadRequest("Invalid URL")
	}
	user, err := u.RepoUser().GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return u.ResponseNotFound("User not found with this email")
		}
		return u.ResponseInternalServer("Failed to update password")
	}

	hashedPassword, err := auth.GenerateHashedPassword(password)
	if err != nil {
		return u.ResponseInternalServer("Failed to update password")
	}
	err = u.RepoUser().UpdatePasswordWithEmail(ctx, email, string(hashedPassword))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return u.ResponseNotFound("User not found with this email")
		}
		return u.ResponseInternalServer("Failed to update the password")
	}
	if err := u.Cache().GetRefreshTokenSession().InvalidateRefreshToken(user.ID); err != nil {
		log.Default().Println("Failed to invalidate refresh token after password update:", err)
	}

	return u.ResponseOK("Password updated successfully", nil)
}
