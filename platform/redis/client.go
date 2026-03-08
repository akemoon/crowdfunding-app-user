package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func NewClient(ctx context.Context, url string) (*redis.Client, error) {
	options, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(options)

	// TODO: maybe use ctx background

	err = client.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}

	return client, nil
}
