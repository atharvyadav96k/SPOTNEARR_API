package services

import (
	"context"
	"log"

	"github.com/atharvyadav96k/SPOTNEARR_API/auth"
	"github.com/atharvyadav96k/SPOTNEARR_API/dtos"
	"github.com/atharvyadav96k/SPOTNEARR_API/models"
	"github.com/atharvyadav96k/SPOTNEARR_API/repository"
	"github.com/atharvyadav96k/SPOTNEARR_API/utils/response"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	base_service
	repo repository.IUserRepository
}

func NewUserService(repo repository.IUserRepository) *UserService {
	return &UserService{repo: repo}
}

func (u *UserService) RegisterUser(user *models.User) response.Res {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(*user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return u.ResponseBadRequest("Failed to process password")
	}
	hashedPassword := string(hashedBytes)
	user.PasswordHash = &hashedPassword
	if err := u.repo.Register(context.Background(), user); err != nil {
		return u.ResponseBadRequest(err.Error())
	}
	return u.ResponseCreated("User registered successfully", nil)
}

func (u *UserService) Login(email string, password string) response.Res {
	user, err := u.repo.GetByEmail(context.Background(), email)
	if err != nil {
		return u.ResponseBadRequest("Invalid email or password")
	}
	err = bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(password))
	if err != nil {
		log.Default().Println(err)
		return u.ResponseBadRequest("Invalid email or password")
	}
	refreshToken, err := auth.GenerateRefreshToken(user.ID, user.BusinessID, "dummy")
	if err != nil {
		log.Default().Println(err)
		return u.ResponseInternalServer("Failed to login")
	}
	err = u.repo.SetRefreshToken(context.Background(), user.ID, refreshToken)
	if err != nil {
		log.Default().Println(err)
		return u.ResponseInternalServer("Failed to login")
	}
	accessToken, err := auth.GenerateAccessToken(user.ID, user.BusinessID, "dummy")
	if err != nil {
		log.Default().Println(err)
		return u.ResponseInternalServer("Failed to login")
	}
	var data dtos.User
	data = data.ResponseMapper(user, accessToken, refreshToken)
	return u.ResponseOK("Logged in successfully", data)
}

func (u *UserService) Refresh(claims auth.UserClaims, refreshToken string) response.Res {
	user, err := u.repo.GetById(context.Background(), claims.UserId)
	if err != nil {
		log.Default().Println("Failed to get user from db")
		return u.ResponseUnauthorized()
	}
	if user.RefreshToken != refreshToken {
		log.Default().Println(user.RefreshToken)
		log.Default().Println(refreshToken)
		log.Default().Println("Failed to match refresh token from database")
		return u.ResponseUnauthorized()
	}
	accessToken, err := auth.GenerateAccessToken(user.ID, user.BusinessID, "dummy")
	if err != nil {
		log.Default().Println("Failed to generate access token")
		return u.ResponseUnauthorized()
	}
	var data dtos.User
	data = data.ResponseMapper(user, accessToken, refreshToken)
	return u.ResponseOK("New access token", data)
}
