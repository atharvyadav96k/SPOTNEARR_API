package services

import (
	"context"
	"errors"
	"log"

	"github.com/atharvyadav96k/SPOTNEARR_API/auth"
	"github.com/atharvyadav96k/SPOTNEARR_API/connections/cache"
	"github.com/atharvyadav96k/SPOTNEARR_API/models"
	"github.com/atharvyadav96k/SPOTNEARR_API/utils/response"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	base_service
}

func NewUserService(db *gorm.DB, cache *cache.Cache) *UserService {
	return &UserService{
		base_service: NewBaseService(db, cache),
	}
}

func (u *UserService) RegisterUser(user *models.User) response.Res {
	if user == nil {
		return u.ResponseBadRequest("Failed to get user")
	}
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(*user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return u.ResponseBadRequest("Failed to process password")
	}
	hashedPassword := string(hashedBytes)
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

	var currentBusinessID *uint
	var userRole = models.RoleUser
	acccess, err := u.RepoAccess().GetAccessByUserId(ctx, user.ID)
	if err == nil {
		currentBusinessID = &acccess.BusinessID
		userRole = acccess.Role
	}
	var refreshToken string

	if _, err := auth.ValidateToken(user.RefreshToken, "dummy", auth.TypeRefreshToken); err == nil {
		refreshToken = user.RefreshToken
	} else {
		refreshToken, err = auth.GenerateRefreshToken(user.ID, currentBusinessID, "dummy", userRole)
		if err != nil {
			log.Default().Println("Failed to generate refresh token:", err)
			return u.ResponseInternalServer("Failed to login")
		}

		err = u.RepoUser().SetRefreshToken(ctx, user.ID, refreshToken)
		if err != nil {
			log.Default().Println("Failed to save refresh token:", err)
			return u.ResponseInternalServer("Failed to login")
		}
	}
	accessToken, err := auth.GenerateAccessToken(user.ID, currentBusinessID, "dummy", userRole)
	if err != nil {
		log.Default().Println("Failed to generate access token:", err)
		return u.ResponseInternalServer("Failed to login")
	}

	return u.ResponseOK("Logged in successfully", auth.TokenResponse{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	})
}

func (u *UserService) Refresh(claims auth.UserClaims, refreshToken string) response.Res {
	ctx := context.Background()
	user, err := u.RepoUser().GetById(ctx, claims.UserId)
	if err != nil {
		log.Default().Println("Failed to get user from db")
		return u.ResponseUnauthorized()
	}
	if user.RefreshToken != refreshToken {
		return u.ResponseUnauthorized()
	}

	var currentBusinessID *uint
	var userRole = models.RoleUser
	userAccess, err := u.RepoAccess().GetAccessByUserId(ctx, user.ID)
	if err == nil {
		currentBusinessID = &userAccess.BusinessID
		userRole = userAccess.Role
	}

	isOutOfSync := false
	if claims.BusinessId == nil && currentBusinessID != nil {
		isOutOfSync = true
	} else if claims.BusinessId != nil && currentBusinessID == nil {
		isOutOfSync = true
	} else if claims.BusinessId != nil && currentBusinessID != nil && *claims.BusinessId != *currentBusinessID {
		isOutOfSync = true
	} else if claims.UserRole != userRole {
		isOutOfSync = true
	}
	if isOutOfSync {
		newRefreshToken, err := auth.GenerateRefreshToken(user.ID, currentBusinessID, "dummy", userRole)
		if err == nil {
			refreshToken = newRefreshToken
			_ = u.RepoUser().SetRefreshToken(ctx, user.ID, newRefreshToken)
		} else {
			log.Default().Println("Warning: Failed to generate updated refresh token:", err)
		}
	}
	accessToken, err := auth.GenerateAccessToken(user.ID, currentBusinessID, "dummy", userRole)
	if err != nil {
		log.Default().Println("Failed to generate access token")
		return u.ResponseUnauthorized()
	}

	return u.ResponseOK("New access token", auth.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (u *UserService) DismissRefreshToken(id uint) response.Res {
	if err := u.RepoUser().RemoveRefreshToken(context.Background(), id); err != nil {
		return u.ResponseNotFound("Account not found")
	}
	return u.ResponseOK("successfully logout from all devices", nil)
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

func (u *UserService) ResetSessionNotification(email string) response.Res {
	user, err := u.RepoUser().GetByEmail(context.Background(), email)
	if err != nil {
		return u.ResponseInternalServer(err.Error())
	}
	if user.Email == nil {
		return u.ResponseNotFound("User not found")
	}
	return u.ResponseOK("", nil)
}

func (u *UserService) UpdatePassword(userId uint, password string) response.Res {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return u.ResponseInternalServer("Failed to update password")
	}
	err = u.RepoUser().UpdatePasswordWithUserId(context.Background(), userId, string(hashedPassword))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return u.ResponseConflict("No user found")
		}
		return u.ResponseBadRequest("Failed to update password")
	}
	return u.ResponseOK("successfully updated user password", nil)
}
