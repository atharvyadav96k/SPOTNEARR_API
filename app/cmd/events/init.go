package events

import (
	"context"

	"github.com/atharvyadav96k/SPOTNEARR_API/connections/cache"
	"github.com/atharvyadav96k/spotnearr/pkg/mq"
	"gorm.io/gorm"
)

// InitConsumers wires up all MQ consumers for the user service.
func InitConsumers(ctx context.Context, sub *mq.Subscriber, db *gorm.DB, c *cache.Cache) error {
	return NewDealConsumer(sub, db, c).Run(ctx)
}
