package events

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/atharvyadav96k/spotnearr/vendor-svc/models"
	"gorm.io/gorm"
)

// Flusher reads unprocessed search_sync_outboxes rows from the vendor DB and
// pushes each payload to the search service. Vendor owns its outbox lifecycle;
// search never touches the vendor database.
type Flusher struct {
	db     *gorm.DB
	url    string
	client *http.Client
}

func NewFlusher(db *gorm.DB, searchServiceURL string) *Flusher {
	return &Flusher{
		db:     db,
		url:    searchServiceURL + "/internal/sync",
		client: &http.Client{Timeout: 5 * time.Second},
	}
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
		if err := f.push(ctx, row.Payload); err != nil {
			log.Printf("flusher: push row %d: %v", row.ID, err)
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

func (f *Flusher) push(ctx context.Context, payload []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, f.url, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := f.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("search returned %d", resp.StatusCode)
	}
	return nil
}
