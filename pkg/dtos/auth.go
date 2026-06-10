package dtos

import (
	"fmt"
	"strings"

	"github.com/atharvyadav96k/spotnearr/pkg/validate"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

func (r *RegisterRequest) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return fmt.Errorf("name is required")
	}
	if !validate.Email(r.Email) {
		return fmt.Errorf("invalid email address")
	}
	if !validate.Phone(r.Phone) {
		return fmt.Errorf("invalid phone number")
	}
	ok, msg := validate.Password(r.Password)
	if !ok {
		return fmt.Errorf("%s", msg)
	}
	return nil
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (l *LoginRequest) Validate() error {
	if strings.TrimSpace(l.Email) == "" {
		return fmt.Errorf("email is required")
	}
	if strings.TrimSpace(l.Password) == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}

type SessionRequest struct {
	Email string `json:"email"`
}

func (s *SessionRequest) Validate() error {
	if !validate.Email(s.Email) {
		return fmt.Errorf("invalid email address")
	}
	return nil
}

type ResetPasswordRequest struct {
	Password string `json:"password"`
}

func (r *ResetPasswordRequest) Validate() error {
	ok, msg := validate.Password(r.Password)
	if !ok {
		return fmt.Errorf("%s", msg)
	}
	return nil
}
