package user

import (
	"context"

	lib "github.com/akemoon/crowdfunding-app-user/lib/domain"
	"github.com/akemoon/crowdfunding-app-user/modules/user/domain"
	"github.com/google/uuid"
)

type Repo interface {
	CreateUser(ctx context.Context, req lib.CreateUserReq) (uuid.UUID, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, req domain.UpdateProfileReq) (domain.User, error)
	Follow(ctx context.Context, followerID uuid.UUID, followeeID uuid.UUID) error
	Unfollow(ctx context.Context, followerID uuid.UUID, followeeID uuid.UUID) error
	SearchUsers(ctx context.Context, req domain.SearchUsersReq) ([]domain.User, error)
	GetSubscriptions(ctx context.Context, userID uuid.UUID) ([]domain.User, error)
}
