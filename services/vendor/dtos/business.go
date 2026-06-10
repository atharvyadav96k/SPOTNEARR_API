package dtos

import (
	"fmt"
	"strings"
)

type RegisterBusiness struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Desc  string `json:"desc"`
}

func (r *RegisterBusiness) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return fmt.Errorf("name is required")
	}
	if strings.TrimSpace(r.Email) == "" {
		return fmt.Errorf("email is required")
	}
	if strings.TrimSpace(r.Phone) == "" {
		return fmt.Errorf("phone is required")
	}
	return nil
}

type UpdateBusiness struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

func (u *UpdateBusiness) Validate() error {
	if strings.TrimSpace(u.Name) == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}
