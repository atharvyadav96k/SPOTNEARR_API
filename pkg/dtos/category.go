package dtos

import (
	"fmt"
	"strings"
)

type CategoryAddRequest struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func (c *CategoryAddRequest) Validate() error {
	if strings.TrimSpace(c.Name) == "" {
		return fmt.Errorf("name is required")
	}
	if strings.TrimSpace(c.Slug) == "" {
		return fmt.Errorf("slug is required")
	}
	return nil
}
