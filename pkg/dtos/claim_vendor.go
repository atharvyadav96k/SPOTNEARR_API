package dtos

import "fmt"

// ClaimStatusUpdateRequest is sent by the vendor service to update a claim status.
type ClaimStatusUpdateRequest struct {
	Status string `json:"status"` // "accepted" | "rejected"
}

func (c *ClaimStatusUpdateRequest) Validate() error {
	if c.Status != "accepted" && c.Status != "rejected" {
		return fmt.Errorf("status must be 'accepted' or 'rejected'")
	}
	return nil
}
