package events

import (
	"context"
	"encoding/json"
	"log"

	"github.com/atharvyadav96k/spotnearr/pkg/mq"
	vendormodel "github.com/Developer-Aadesh/spotnearr-database/vendordb"
)

// FollowConsumer listens for follow/unfollow events and keeps the
// business follower_count column in sync.
type FollowConsumer struct {
	sub     *mq.Subscriber
	bizRepo vendormodel.IBusinessesRepository
}

func NewFollowConsumer(sub *mq.Subscriber, bizRepo vendormodel.IBusinessesRepository) *FollowConsumer {
	return &FollowConsumer{sub: sub, bizRepo: bizRepo}
}

// Run subscribes to both follow and unfollow topics. Blocks until ctx is done.
func (f *FollowConsumer) Run(ctx context.Context) error {
	if err := f.sub.Subscribe(ctx, "vendor.business.follow", mq.TopicBusinessFollow, f.handleFollow); err != nil {
		return err
	}
	return f.sub.Subscribe(ctx, "vendor.business.unfollow", mq.TopicBusinessUnfollow, f.handleUnfollow)
}

func (f *FollowConsumer) handleFollow(ctx context.Context, body []byte) error {
	var p mq.BusinessFollowPayload
	if err := json.Unmarshal(body, &p); err != nil {
		return err
	}
	if err := f.bizRepo.IncrementFollowerCount(ctx, p.BusinessID); err != nil {
		log.Printf("follow_consumer: increment biz %d: %v", p.BusinessID, err)
		return err
	}
	return nil
}

func (f *FollowConsumer) handleUnfollow(ctx context.Context, body []byte) error {
	var p mq.BusinessFollowPayload
	if err := json.Unmarshal(body, &p); err != nil {
		return err
	}
	if err := f.bizRepo.DecrementFollowerCount(ctx, p.BusinessID); err != nil {
		log.Printf("follow_consumer: decrement biz %d: %v", p.BusinessID, err)
		return err
	}
	return nil
}
