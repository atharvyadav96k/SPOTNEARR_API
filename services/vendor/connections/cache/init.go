package cache

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func InitCache(cacheUrl string, password string) (*Cache, error) {
	opt, err := redis.ParseURL(cacheUrl)
	if err != nil {
		return nil, err
	}
	if password != "" {
		opt.Password = password
	}
	client := redis.NewClient(opt)
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	log.Default().Println("Vendor Service: Redis connected.")
	return newCache(client), nil
}
