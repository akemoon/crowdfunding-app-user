package user

import (
	"context"

	"github.com/akemoon/crowdfunding-app-user/modules/user/domain"
	"github.com/google/uuid"
)

type Repo interface {
	CreateUser(ctx context.Context, req domain.CreateUserReq, passwordHash string) (uuid.UUID, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	ListUsers(ctx context.Context) ([]domain.User, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, req domain.UpdateProfileReq) (domain.User, error)
	Follow(ctx context.Context, followerID uuid.UUID, followeeID uuid.UUID) error
	Unfollow(ctx context.Context, followerID uuid.UUID, followeeID uuid.UUID) error
	GetFollows(ctx context.Context, userID uuid.UUID) ([]domain.User, error)
}
