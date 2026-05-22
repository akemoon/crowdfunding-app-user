package token

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Repo interface {
	Set(ctx context.Context, hash string, userID uuid.UUID, ttl time.Duration) error
	Get(ctx context.Context, hash string) (uuid.UUID, error)
	Delete(ctx context.Context, hash string) error
	DeleteAllByUserID(ctx context.Context, userID uuid.UUID) error
}
