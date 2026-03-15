package redis

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RefreshTokenRepo struct {
	redisClient *redis.Client
}

func NewRefreshTokenRepo(rc *redis.Client) *RefreshTokenRepo {
	return &RefreshTokenRepo{
		redisClient: rc,
	}
}

func (r *RefreshTokenRepo) Set(ctx context.Context, hash string, userID uuid.UUID, ttl time.Duration) error {
	return r.redisClient.Set(ctx, hash, userID.String(), ttl).Err()
}

func (r *RefreshTokenRepo) Get(ctx context.Context, hash string) (uuid.UUID, error) {
	val, err := r.redisClient.Get(ctx, hash).Result()
	if err != nil {
		return uuid.UUID{}, err
	}

	return uuid.Parse(val)
}

func (r *RefreshTokenRepo) Delete(ctx context.Context, hash string) error {
	return r.redisClient.Del(ctx, hash).Err()
}
