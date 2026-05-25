package dtos

import (
	"github.com/atharvyadav96k/SPOTNEARR_API/models"
)

type User struct {
	FullName     string `json:"full_name,omitempty"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Password     string `json:"password,omitempty"`
	Role         string `json:"role"`
	AccessToken  string `json:"access_token",omitempty`
	RefreshToken string `json:"refresh_token"`
}

func (u *User) RequestMapper(data User) *models.User {
	emailStr := data.Email
	phoneStr := data.Phone
	passStr := data.Password

	return &models.User{
		FullName:     data.FullName,
		Email:        &emailStr,
		Phone:        &phoneStr,
		Role:         models.UserRole(data.Role),
		PasswordHash: &passStr,
	}
}

func (u *User) ResponseMapper(data *models.User, access_token string, refresh_token string) User {
	var email, phone string

	if data.Email != nil {
		email = *data.Email
	}
	if data.Phone != nil {
		phone = *data.Phone
	}

	return User{
		FullName:     data.FullName,
		Email:        email,
		Phone:        phone,
		Role:         string(data.Role),
		AccessToken:  access_token,
		RefreshToken: refresh_token,
	}
}
