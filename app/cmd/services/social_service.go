package services

import (
	"context"
	"errors"
	"log"

	"github.com/atharvyadav96k/SPOTNEARR_API/utils/response"
	"github.com/atharvyadav96k/spotnearr/pkg/mq"
	"gorm.io/gorm"
)

type SocialService struct {
	base_service
	publisher *mq.Publisher
}

func NewSocialService(b base_service, publisher *mq.Publisher) *SocialService {
	return &SocialService{base_service: b, publisher: publisher}
}

// FollowBusiness records a follow in the user DB then publishes a follow event.
func (s *SocialService) FollowBusiness(userID, bizID uint) response.Res {
	ctx := context.Background()

	already, err := s.RepoFollow().IsFollowing(ctx, userID, bizID)
	if err != nil {
		return s.ResponseInternalServer("Failed to check follow status")
	}
	if already {
		return s.ResponseConflict("Already following this business")
	}

	if err := s.RepoFollow().Follow(ctx, userID, bizID); err != nil {
		return s.ResponseInternalServer("Failed to follow business")
	}

	if err := s.publisher.Publish(ctx, mq.TopicBusinessFollow, mq.BusinessFollowPayload{BusinessID: bizID, UserID: userID}); err != nil {
		log.Printf("social: publish follow event biz %d: %v", bizID, err)
	}
	return s.ResponseOK("Followed successfully", nil)
}

// UnfollowBusiness removes the follow from the user DB then publishes an unfollow event.
func (s *SocialService) UnfollowBusiness(userID, bizID uint) response.Res {
	ctx := context.Background()

	if err := s.RepoFollow().Unfollow(ctx, userID, bizID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return s.ResponseNotFound("Not following this business")
		}
		return s.ResponseInternalServer("Failed to unfollow business")
	}

	if err := s.publisher.Publish(ctx, mq.TopicBusinessUnfollow, mq.BusinessFollowPayload{BusinessID: bizID, UserID: userID}); err != nil {
		log.Printf("social: publish unfollow event biz %d: %v", bizID, err)
	}
	return s.ResponseOK("Unfollowed successfully", nil)
}

// LikeProduct records a product like for the user.
func (s *SocialService) LikeProduct(userID, invProductID uint) response.Res {
	if err := s.RepoProductEngage().Like(context.Background(), userID, invProductID); err != nil {
		return s.ResponseInternalServer("Failed to like product")
	}
	return s.ResponseOK("Product liked", nil)
}

// UnlikeProduct removes a product like for the user.
func (s *SocialService) UnlikeProduct(userID, invProductID uint) response.Res {
	if err := s.RepoProductEngage().Unlike(context.Background(), userID, invProductID); err != nil {
		return s.ResponseInternalServer("Failed to unlike product")
	}
	return s.ResponseOK("Product unliked", nil)
}

// SaveProduct saves a product for the user.
func (s *SocialService) SaveProduct(userID, invProductID uint) response.Res {
	if err := s.RepoProductEngage().Save(context.Background(), userID, invProductID); err != nil {
		return s.ResponseInternalServer("Failed to save product")
	}
	return s.ResponseOK("Product saved", nil)
}

// UnsaveProduct removes a saved product for the user.
func (s *SocialService) UnsaveProduct(userID, invProductID uint) response.Res {
	if err := s.RepoProductEngage().Unsave(context.Background(), userID, invProductID); err != nil {
		return s.ResponseInternalServer("Failed to unsave product")
	}
	return s.ResponseOK("Product unsaved", nil)
}

// LikeSpotlight records a spotlight like for the user.
func (s *SocialService) LikeSpotlight(userID, spotlightID uint) response.Res {
	if err := s.RepoSpotlightEngage().Like(context.Background(), userID, spotlightID); err != nil {
		return s.ResponseInternalServer("Failed to like spotlight")
	}
	return s.ResponseOK("Spotlight liked", nil)
}

// UnlikeSpotlight removes a spotlight like for the user.
func (s *SocialService) UnlikeSpotlight(userID, spotlightID uint) response.Res {
	if err := s.RepoSpotlightEngage().Unlike(context.Background(), userID, spotlightID); err != nil {
		return s.ResponseInternalServer("Failed to unlike spotlight")
	}
	return s.ResponseOK("Spotlight unliked", nil)
}

// SaveSpotlight saves a spotlight for the user.
func (s *SocialService) SaveSpotlight(userID, spotlightID uint) response.Res {
	if err := s.RepoSpotlightEngage().Save(context.Background(), userID, spotlightID); err != nil {
		return s.ResponseInternalServer("Failed to save spotlight")
	}
	return s.ResponseOK("Spotlight saved", nil)
}

// UnsaveSpotlight removes a saved spotlight for the user.
func (s *SocialService) UnsaveSpotlight(userID, spotlightID uint) response.Res {
	if err := s.RepoSpotlightEngage().Unsave(context.Background(), userID, spotlightID); err != nil {
		return s.ResponseInternalServer("Failed to unsave spotlight")
	}
	return s.ResponseOK("Spotlight unsaved", nil)
}
