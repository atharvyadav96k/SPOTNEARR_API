package dtos

import "fmt"

type UpdateProfileRequest struct {
	FullName string `json:"fullName"`
}

func (r UpdateProfileRequest) Validate() error {
	if r.FullName == "" {
		return fmt.Errorf("fullName is required")
	}
	return nil
}
