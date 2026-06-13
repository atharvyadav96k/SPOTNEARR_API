package dtos

import (
	"fmt"
	"strings"
	"time"
)

type OfferCreateRequest struct {
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	DiscountType  string     `json:"discountType"` // "flat" | "percent"
	DiscountValue float64    `json:"discountValue"`
	MinOrderValue *float64   `json:"minOrderValue"`
	Code          *string    `json:"code"`
	MaxUsage      *int       `json:"maxUsage"`
	ExpiresAt     *time.Time `json:"expiresAt"`
}

func (o *OfferCreateRequest) Validate() error {
	if strings.TrimSpace(o.Title) == "" {
		return fmt.Errorf("title is required")
	}
	if o.DiscountType != "flat" && o.DiscountType != "percent" {
		return fmt.Errorf("discountType must be 'flat' or 'percent'")
	}
	if o.DiscountValue <= 0 {
		return fmt.Errorf("discountValue must be greater than 0")
	}
	if o.DiscountType == "percent" && o.DiscountValue > 100 {
		return fmt.Errorf("percent discount cannot exceed 100")
	}
	return nil
}

type OfferUpdateRequest struct {
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	DiscountType  string     `json:"discountType"`
	DiscountValue float64    `json:"discountValue"`
	MinOrderValue *float64   `json:"minOrderValue"`
	Code          *string    `json:"code"`
	MaxUsage      *int       `json:"maxUsage"`
	Active        *bool      `json:"active"`
	ExpiresAt     *time.Time `json:"expiresAt"`
}

func (o *OfferUpdateRequest) Validate() error {
	if o.DiscountType != "" && o.DiscountType != "flat" && o.DiscountType != "percent" {
		return fmt.Errorf("discountType must be 'flat' or 'percent'")
	}
	if o.DiscountValue < 0 {
		return fmt.Errorf("discountValue cannot be negative")
	}
	return nil
}
