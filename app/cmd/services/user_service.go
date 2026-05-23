package services

import (
	"context"

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
		return u.ResponseBadRequest("Invalid email or password")
	}
	var data dtos.User
	data = data.ResponseMapper(user)
	return u.ResponseOK("Logged in successfully", data)
}
