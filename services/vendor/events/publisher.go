package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	vendordb "github.com/Developer-Aadesh/spotnearr-database/vendordb"
	"github.com/atharvyadav96k/spotnearr/pkg/mq"
	"github.com/atharvyadav96k/spotnearr/vendor-svc/models"
	"gorm.io/gorm"
)

// Flusher reads unprocessed search_sync_outboxes rows from the vendor DB and
// publishes each payload to the search service via RabbitMQ.
type Flusher struct {
	db  *gorm.DB
	pub *mq.Publisher
}

func NewFlusher(db *gorm.DB, pub *mq.Publisher) *Flusher {
	return &Flusher{db: db, pub: pub}
}

// Run starts the fallback poll loop. Call in a goroutine.
func (f *Flusher) Run(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			f.Flush(ctx)
		}
	}
}

// Flush reads unprocessed outbox rows inside a transaction with SKIP LOCKED so
// concurrent vendor instances don't double-deliver the same rows.
func (f *Flusher) Flush(ctx context.Context) {
	tx := f.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		log.Printf("flusher: begin tx: %v", tx.Error)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var rows []models.SearchSyncOutbox
	if err := tx.Raw(`
		SELECT id, event_type, payload, created_at
		FROM search_sync_outboxes
		WHERE processed = false
		ORDER BY id ASC
		LIMIT 100
		FOR UPDATE SKIP LOCKED
	`).Scan(&rows).Error; err != nil {
		log.Printf("flusher: fetch outbox: %v", err)
		tx.Rollback()
		return
	}

	if len(rows) == 0 {
		tx.Rollback()
		return
	}

	for _, row := range rows {
		if err := f.publish(ctx, row.Payload); err != nil {
			log.Printf("flusher: publish row %d: %v", row.ID, err)
			continue
		}
		now := time.Now()
		if err := tx.Exec(
			`UPDATE search_sync_outboxes SET processed = true, processed_at = ? WHERE id = ?`,
			now, row.ID,
		).Error; err != nil {
			log.Printf("flusher: mark row %d: %v", row.ID, err)
		}
	}

	tx.Commit()
}

func (f *Flusher) publish(ctx context.Context, payload []byte) error {
	var p vendordb.OutboxPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return fmt.Errorf("unmarshal outbox payload: %w", err)
	}
	return f.pub.Publish(ctx, mq.TopicProductSync, p)
}
