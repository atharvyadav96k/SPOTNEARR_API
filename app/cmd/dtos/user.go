package dtos

import (
	"github.com/atharvyadav96k/SPOTNEARR_API/models"
)

type User struct {
	FullName string `json:"full_name,omitempty"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role"`
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

func (u *User) ResponseMapper(data *models.User) User {
	var email, phone string

	if data.Email != nil {
		email = *data.Email
	}
	if data.Phone != nil {
		phone = *data.Phone
	}

	return User{
		FullName: data.FullName,
		Email:    email,
		Phone:    phone,
		Role:     string(data.Role),
	}
}
