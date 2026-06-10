package services

import (
	"context"
	"errors"

	"github.com/atharvyadav96k/SPOTNEARR_API/utils/response"
	usermodel "github.com/Developer-Aadesh/spotnearr-database/user"
	"gorm.io/gorm"
)

type ReviewService struct {
	base_service
}

func NewReviewService(b base_service) *ReviewService {
	return &ReviewService{base_service: b}
}

func (s *ReviewService) AddReview(userID uint, targetType usermodel.ReviewTarget, targetID uint, stars uint8, comment string) response.Res {
	if stars < 1 || stars > 5 {
		return s.ResponseBadRequest("Stars must be between 1 and 5")
	}
	if targetID == 0 {
		return s.ResponseBadRequest("Target ID is required")
	}

	ctx := context.Background()

	_, err := s.RepoReview().GetByUserAndTarget(ctx, userID, targetType, targetID)
	if err == nil {
		return s.ResponseConflict("You have already reviewed this")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return s.ResponseInternalServer("Failed to check existing review")
	}

	review := usermodel.Review{
		UserID:     userID,
		TargetType: targetType,
		TargetID:   targetID,
		Stars:      stars,
		Comment:    comment,
	}
	created, err := s.RepoReview().Add(ctx, review)
	if err != nil {
		return s.ResponseInternalServer("Failed to add review")
	}
	return s.ResponseCreated("Review added successfully", created)
}

func (s *ReviewService) UpdateReview(userID uint, targetType usermodel.ReviewTarget, targetID uint, stars uint8, comment string) response.Res {
	if stars < 1 || stars > 5 {
		return s.ResponseBadRequest("Stars must be between 1 and 5")
	}
	if targetID == 0 {
		return s.ResponseBadRequest("Target ID is required")
	}

	updated, err := s.RepoReview().Update(context.Background(), userID, targetType, targetID, stars, comment)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return s.ResponseNotFound("Review not found")
		}
		return s.ResponseInternalServer("Failed to update review")
	}
	return s.ResponseOK("Review updated successfully", updated)
}

func (s *ReviewService) DeleteReview(userID uint, targetType usermodel.ReviewTarget, targetID uint) response.Res {
	if targetID == 0 {
		return s.ResponseBadRequest("Target ID is required")
	}

	err := s.RepoReview().Delete(context.Background(), userID, targetType, targetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return s.ResponseNotFound("Review not found")
		}
		return s.ResponseInternalServer("Failed to delete review")
	}
	return s.ResponseOK("Review deleted successfully", nil)
}

func (s *ReviewService) GetReviews(targetType usermodel.ReviewTarget, targetID uint) response.Res {
	if targetID == 0 {
		return s.ResponseBadRequest("Target ID is required")
	}

	reviews, err := s.RepoReview().GetByTarget(context.Background(), targetType, targetID)
	if err != nil {
		return s.ResponseInternalServer("Failed to fetch reviews")
	}
	return s.ResponseOK("Reviews fetched successfully", reviews)
}
