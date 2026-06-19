package cache

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const dealGeoCacheTTL = time.Hour

type dealGeoCache struct{ *base_cache }

func newDealGeoCache(client *redis.Client) *dealGeoCache {
	return &dealGeoCache{base_cache: newBaseCache(client)}
}

func (d *dealGeoCache) key(geohash5 string) string {
	return d.pk("deals:geo:" + geohash5)
}

func (d *dealGeoCache) AddDeal(ctx context.Context, geohash5 string, dealID uint) error {
	k := d.key(geohash5)
	pipe := d.getClient().Pipeline()
	pipe.SAdd(ctx, k, fmt.Sprintf("%d", dealID))
	pipe.Expire(ctx, k, dealGeoCacheTTL)
	_, err := pipe.Exec(ctx)
	return err
}

func (d *dealGeoCache) RemoveDeal(ctx context.Context, geohash5 string, dealID uint) error {
	return d.getClient().SRem(ctx, d.key(geohash5), fmt.Sprintf("%d", dealID)).Err()
}

// GetDealIDs returns nil on cache miss (key absent), empty slice if key present but empty.
func (d *dealGeoCache) GetDealIDs(ctx context.Context, geohash5 string) ([]uint, error) {
	k := d.key(geohash5)
	exists, err := d.getClient().Exists(ctx, k).Result()
	if err != nil || exists == 0 {
		return nil, err
	}
	members, err := d.getClient().SMembers(ctx, k).Result()
	if err != nil {
		return nil, err
	}
	ids := make([]uint, 0, len(members))
	for _, m := range members {
		if id, err := strconv.ParseUint(m, 10, 64); err == nil {
			ids = append(ids, uint(id))
		}
	}
	return ids, nil
}

func (d *dealGeoCache) PopulateDealIDs(ctx context.Context, geohash5 string, dealIDs []uint) error {
	k := d.key(geohash5)
	pipe := d.getClient().Pipeline()
	pipe.Del(ctx, k)
	if len(dealIDs) > 0 {
		members := make([]interface{}, len(dealIDs))
		for i, id := range dealIDs {
			members[i] = fmt.Sprintf("%d", id)
		}
		pipe.SAdd(ctx, k, members...)
	}
	pipe.Expire(ctx, k, dealGeoCacheTTL)
	_, err := pipe.Exec(ctx)
	return err
}
