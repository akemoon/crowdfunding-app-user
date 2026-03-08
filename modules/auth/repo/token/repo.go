package token

import (
	"context"
	"time"
)

type Repo interface {
	Set(ctx context.Context, refreshToken string, ttl time.Duration) error
	Check(ctx context.Context, refreshToken string) error
	Delete(ctx context.Context, refreshToken string) error
}
