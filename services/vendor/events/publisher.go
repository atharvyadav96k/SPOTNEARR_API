package events

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

const SearchSyncChannel = "search:sync"

// Publisher sends lightweight notifications to the Search Service after
// outbox rows are committed. The Search Service Redis subscriber picks these
// up immediately; the 30s outbox poller is the fallback when this fails.
type Publisher struct{ client *redis.Client }

func NewPublisher(client *redis.Client) *Publisher {
	return &Publisher{client: client}
}

// Notify publishes a ping on the search:sync channel. Content is ignored by
// the subscriber — it just wakes up the poller immediately.
func (p *Publisher) Notify(ctx context.Context) {
	if err := p.client.Publish(ctx, SearchSyncChannel, "sync").Err(); err != nil {
		log.Printf("events: redis publish search:sync: %v", err)
	}
}
