package freq

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	searchdb "github.com/Developer-Aadesh/spotnearr-database/search"
	searchpostgres "github.com/Developer-Aadesh/spotnearr-database/search/postgres"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const (
	deltaHashKey   = "search:tf_delta"
	flushingHashKey = "search:tf_flushing"
)

// Flusher periodically moves token-category frequency deltas from Redis into the database.
// Writes are first accumulated in a Redis hash (HINCRBY) and then flushed in batches,
// keeping per-write DB pressure near zero.
type Flusher struct {
	repo *searchpostgres.SearchRepository
	rdb  *redis.Client
}

func NewFlusher(db *gorm.DB, rdb *redis.Client) *Flusher {
	return &Flusher{
		repo: searchpostgres.NewSearchRepository(db),
		rdb:  rdb,
	}
}

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

// Flush atomically renames the accumulator hash so new writes land on a fresh key
// while the previous batch is being written to the database.
func (f *Flusher) Flush(ctx context.Context) {
	if err := f.rdb.Rename(ctx, deltaHashKey, flushingHashKey).Err(); err != nil {
		if strings.Contains(err.Error(), "no such key") {
			return // nothing accumulated
		}
		log.Printf("freq flush: rename: %v", err)
		return
	}

	fields, err := f.rdb.HGetAll(ctx, flushingHashKey).Result()
	if err != nil {
		log.Printf("freq flush: hgetall: %v", err)
		f.rdb.Del(ctx, flushingHashKey)
		return
	}

	deltas := make([]searchdb.FreqDelta, 0, len(fields))
	for field, val := range fields {
		colonIdx := strings.LastIndex(field, ":")
		if colonIdx < 1 {
			continue
		}
		token := field[:colonIdx]
		catID, err := strconv.ParseUint(field[colonIdx+1:], 10, 64)
		if err != nil {
			continue
		}
		delta, err := strconv.ParseInt(val, 10, 64)
		if err != nil || delta == 0 {
			continue
		}
		deltas = append(deltas, searchdb.FreqDelta{
			Token:      token,
			CategoryID: uint(catID),
			Delta:      delta,
		})
	}

	if len(deltas) > 0 {
		if err := f.repo.FlushFreqDeltas(ctx, deltas); err != nil {
			log.Printf("freq flush: db write: %v", err)
			// On failure, merge unwritten deltas back so the next tick retries them.
			pipe := f.rdb.Pipeline()
			for _, d := range deltas {
				pipe.HIncrBy(ctx, deltaHashKey, fmt.Sprintf("%s:%d", d.Token, d.CategoryID), d.Delta)
			}
			if _, pErr := pipe.Exec(ctx); pErr != nil {
				log.Printf("freq flush: merge-back: %v", pErr)
			}
		}
	}

	f.rdb.Del(ctx, flushingHashKey)
}
