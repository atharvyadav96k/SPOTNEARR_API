package sync

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

const searchSyncChannel = "search:sync"

// Subscriber listens on the Redis search:sync channel. Each message triggers
// an immediate ProcessPending call so newly added outbox rows are applied within
// milliseconds of the vendor write, rather than waiting for the 30s poll tick.
type Subscriber struct {
	client  *redis.Client
	poller  *Poller
}

func NewSubscriber(client *redis.Client, poller *Poller) *Subscriber {
	return &Subscriber{client: client, poller: poller}
}

// Run blocks until ctx is cancelled. Call in a goroutine.
func (s *Subscriber) Run(ctx context.Context) {
	pubsub := s.client.Subscribe(ctx, searchSyncChannel)
	defer pubsub.Close()

	log.Printf("sync: subscribed to Redis channel %q", searchSyncChannel)
	ch := pubsub.Channel()
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}
			log.Printf("sync: received %q on %s — processing pending rows", msg.Payload, msg.Channel)
			s.poller.ProcessPending(ctx)
		}
	}
}
