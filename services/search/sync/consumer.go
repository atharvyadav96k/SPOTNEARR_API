package sync

import (
	"context"

	"github.com/atharvyadav96k/spotnearr/pkg/mq"
)

// RunConsumer subscribes to the vendor.product.sync topic and feeds each
// message into the Applier. Blocks until ctx is done.
func RunConsumer(ctx context.Context, sub *mq.Subscriber, applier *Applier) error {
	return sub.Subscribe(ctx, "search.product.sync", mq.TopicProductSync, applier.Apply)
}
