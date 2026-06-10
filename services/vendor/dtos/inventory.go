package dtos

import (
	"fmt"
	"strings"
)

type InventoryCreateRequest struct {
	Name          string  `json:"name"`
	StreetAddress string  `json:"streetAddress"`
	Lat           float64 `json:"lat"`
	Long          float64 `json:"long"`
}

func (i *InventoryCreateRequest) Validate() error {
	if strings.TrimSpace(i.Name) == "" {
		return fmt.Errorf("name is required")
	}
	if strings.TrimSpace(i.StreetAddress) == "" {
		return fmt.Errorf("streetAddress is required")
	}
	return nil
}

type InventoryUpdateRequest struct {
	Name    string   `json:"name"`
	Address string   `json:"address"`
	Lat     *float64 `json:"lat"`
	Long    *float64 `json:"long"`
}

func (i *InventoryUpdateRequest) Validate() error {
	if i.Lat == nil || i.Long == nil {
		return fmt.Errorf("lat and long are required")
	}
	return nil
}

type InventoryAddProductRequest struct {
	ProductID uint  `json:"productID"`
	Count     *int  `json:"count"`
	Available *bool `json:"available"`
}

func (i *InventoryAddProductRequest) Validate() error {
	if i.ProductID == 0 {
		return fmt.Errorf("productID is required")
	}
	return nil
}

type InventoryUpdateProductRequest struct {
	Count     *int  `json:"count"`
	Available *bool `json:"available"`
}

func (i *InventoryUpdateProductRequest) Validate() error { return nil }
