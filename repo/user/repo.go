package user

import (
	"context"

	"github.com/google/uuid"

	"github.com/akemoon/crowdfunding-app-user/domain"
)

type Repo interface {
	CreateUser(ctx context.Context, req domain.CreateUserReq) (domain.CreateUserResp, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	GetUserByUsername(ctx context.Context, username string) (domain.User, error)
}
