package dtos

import "fmt"

type ReviewRequest struct {
	TargetID uint   `json:"targetId"`
	Stars    uint8  `json:"stars"`
	Comment  string `json:"comment"`
}

func (r *ReviewRequest) Validate() error {
	if r.TargetID == 0 {
		return fmt.Errorf("targetId is required")
	}
	if r.Stars < 1 || r.Stars > 5 {
		return fmt.Errorf("stars must be between 1 and 5")
	}
	return nil
}
