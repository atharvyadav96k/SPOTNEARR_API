package cache

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func InitCache(cacheUrl string, password string) (*Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cacheUrl,
		Password: password,
		DB:       0,
	})

	cache := client.Ping(context.Background())
	if cache.Err() != nil {
		return nil, cache.Err()
	}
	log.Default().Println(cache.Val())
	return newCache(client), nil
}
