package cache

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func InitCache(cacheUrl string) (*Cache, error) {
	opt, err := redis.ParseURL(cacheUrl)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)

	err = client.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}

	log.Default().Println("Connected to local Redis successfully!")
	return newCache(client), nil
}
