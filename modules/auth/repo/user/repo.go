package user

import (
	"context"

	"github.com/akemoon/crowdfunding-app-user/lib/domain"
	"github.com/google/uuid"
)

type Repo interface {
	CreateUser(ctx context.Context, req domain.CreateUserReq) (uuid.UUID, error)
	GetCredentialsByEmail(ctx context.Context, email string) (domain.UserCredentials, error)
	GetCredentialsByID(ctx context.Context, id uuid.UUID) (domain.UserCredentials, error)
}
