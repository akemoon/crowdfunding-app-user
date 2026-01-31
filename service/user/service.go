package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/akemoon/crowdfunding-app-user/domain"
	"github.com/akemoon/crowdfunding-app-user/repo/user"
)

type Service struct {
	repo user.Repo
}

func NewService(repo user.Repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateUser(ctx context.Context, req domain.CreateUserReq) (domain.CreateUserResp, error) {
	err := domain.ValidateUsernameLen(req.Username)
	if err != nil {
		return domain.CreateUserResp{}, fmt.Errorf("%w: %s", domain.ErrInvalidUsername, err)
	}

	resp, err := s.repo.CreateUser(ctx, req)
	if err != nil {
		return domain.CreateUserResp{}, fmt.Errorf("repo: %w", err)
	}

	return resp, nil
}

func (s *Service) GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("repo: %w", err)
	}

	return user, nil
}

func (s *Service) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return domain.User{}, fmt.Errorf("repo: %w", err)
	}

	return user, nil
}
