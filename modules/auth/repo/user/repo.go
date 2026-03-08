package user

import (
	"context"

	"github.com/akemoon/crowdfunding-app-user/lib/domain"
)

type Repo interface {
	CreateUser(ctx context.Context, req domain.CreateUserReq) error
	GetCredentialsByEmail(ctx context.Context, email string) (domain.UserCredentials, error)
}
