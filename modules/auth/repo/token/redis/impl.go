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
	setKey := r.userSetKey(userID.String())
	pipe := r.redisClient.Pipeline()
	pipe.Set(ctx, hash, userID.String(), ttl)
	pipe.SAdd(ctx, setKey, hash)
	pipe.Expire(ctx, setKey, ttl)
	_, err := pipe.Exec(ctx)
	return err
}

func (r *RefreshTokenRepo) Get(ctx context.Context, hash string) (uuid.UUID, error) {
	val, err := r.redisClient.Get(ctx, hash).Result()
	if err != nil {
		return uuid.UUID{}, err
	}

	return uuid.Parse(val)
}

func (r *RefreshTokenRepo) Delete(ctx context.Context, hash string) error {
	userIDStr, err := r.redisClient.Get(ctx, hash).Result()
	if err != nil {
		return r.redisClient.Del(ctx, hash).Err()
	}

	pipe := r.redisClient.Pipeline()
	pipe.Del(ctx, hash)
	pipe.SRem(ctx, r.userSetKey(userIDStr), hash)
	_, err = pipe.Exec(ctx)
	return err
}

func (r *RefreshTokenRepo) DeleteAllByUserID(ctx context.Context, userID uuid.UUID) error {
	setKey := r.userSetKey(userID.String())

	hashes, err := r.redisClient.SMembers(ctx, setKey).Result()
	if err != nil {
		return err
	}

	keys := append(hashes, setKey)
	return r.redisClient.Del(ctx, keys...).Err()
}

func (r *RefreshTokenRepo) userSetKey(userID string) string {
	return "user:" + userID + ":tokens"
}
