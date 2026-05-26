package dtos

import (
	"github.com/atharvyadav96k/SPOTNEARR_API/models"
)

type User struct {
	FullName string `json:"full_name,omitempty"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password,omitempty"`
}

func (u *User) RequestMapper() *models.User {
	emailStr := u.Email
	phoneStr := u.Phone
	passStr := u.Password

	return &models.User{
		FullName:     u.FullName,
		Email:        &emailStr,
		Phone:        &phoneStr,
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
	}
}
