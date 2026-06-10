package dtos

import (
	"fmt"
	"strings"
)

type ValueUnit struct {
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}

type ProductAddRequest struct {
	Name        string    `json:"name"`
	Price       ValueUnit `json:"price"`
	Quantity    ValueUnit `json:"quantity"`
	Desc        string    `json:"desc"`
	CategoryIDs []uint    `json:"categoryIds"`
	StoreIDs    []uint    `json:"storeIds"`
}

func (p *ProductAddRequest) Validate() error {
	if strings.TrimSpace(p.Name) == "" {
		return fmt.Errorf("product name is required")
	}
	if p.Price.Value <= 0 || strings.TrimSpace(p.Price.Unit) == "" {
		return fmt.Errorf("valid price with unit is required")
	}
	if strings.TrimSpace(p.Quantity.Unit) == "" {
		return fmt.Errorf("quantity unit is required")
	}
	if len(p.CategoryIDs) == 0 {
		return fmt.Errorf("at least one category is required")
	}
	return nil
}

type ProductUpdateRequest struct {
	Name     string    `json:"name"`
	Price    ValueUnit `json:"price"`
	Quantity ValueUnit `json:"quantity"`
	Desc     string    `json:"desc"`
}

func (p *ProductUpdateRequest) Validate() error {
	if p.Price.Value <= 0 || strings.TrimSpace(p.Price.Unit) == "" {
		return fmt.Errorf("valid price with unit is required")
	}
	if strings.TrimSpace(p.Quantity.Unit) == "" {
		return fmt.Errorf("quantity unit is required")
	}
	return nil
}
