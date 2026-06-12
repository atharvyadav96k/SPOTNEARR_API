package events

import (
	"context"

	"github.com/atharvyadav96k/spotnearr/pkg/mq"
	vendorpostgres "github.com/Developer-Aadesh/spotnearr-database/vendordb/postgres"
	"gorm.io/gorm"
)

// InitConsumers wires up all MQ consumers for the vendor service.
// Add new consumers here as the service grows.
func InitConsumers(ctx context.Context, sub *mq.Subscriber, db *gorm.DB) error {
	follow := NewFollowConsumer(sub, vendorpostgres.NewBusinessRepository(db))
	return follow.Run(ctx)
}
