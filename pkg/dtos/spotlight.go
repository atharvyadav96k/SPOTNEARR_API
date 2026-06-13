package dtos

import "fmt"

type SpotlightCreateRequest struct {
	Type      string `json:"type"`      // "product" | "offer" | "general"
	Title     string `json:"title"`
	Caption   string `json:"caption"`
	MediaURL  string `json:"mediaUrl"`
	MediaType string `json:"mediaType"` // "image" | "video"
	ProductID *uint  `json:"productId"`
	OfferID   *uint  `json:"offerId"`
}

func (s *SpotlightCreateRequest) Validate() error {
	switch s.Type {
	case "product", "offer", "general":
	default:
		return fmt.Errorf("type must be 'product', 'offer', or 'general'")
	}
	if s.Title == "" {
		return fmt.Errorf("title is required")
	}
	if s.MediaURL == "" {
		return fmt.Errorf("mediaUrl is required")
	}
	if s.MediaType != "image" && s.MediaType != "video" {
		return fmt.Errorf("mediaType must be 'image' or 'video'")
	}
	return nil
}
