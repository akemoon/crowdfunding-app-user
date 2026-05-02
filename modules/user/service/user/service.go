package user

import (
	"context"
	"fmt"

	lib "github.com/akemoon/crowdfunding-app-user/lib/domain"
	"github.com/akemoon/golib/validation"
	"github.com/akemoon/crowdfunding-app-user/modules/user/domain"
	"github.com/akemoon/crowdfunding-app-user/modules/user/repo/user"
	"github.com/google/uuid"
)

type Service struct {
	repo           user.Repo
	avatarsBaseURL string
}

func NewService(repo user.Repo, avatarsBaseURL string) *Service {
	return &Service{
		repo:           repo,
		avatarsBaseURL: avatarsBaseURL,
	}
}

func (s *Service) GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("repo: %w", err)
	}

	user.AvatarUrl = s.avatarsBaseURL + "/" + user.AvatarUrl

	return user, nil
}

func (s *Service) UpdateProfile(ctx context.Context, userID uuid.UUID, req domain.UpdateProfileReq) (domain.User, error) {
	ve := &validation.Error{}

	// TODO: maybe check in cycle
	err := lib.ValidateUsername(req.Username)
	if err != nil {
		ve.Add("username", err.Error())
	}
	err = domain.ValidateDisplayName(req.DisplayName)
	if err != nil {
		ve.Add("displayName", err.Error())
	}
	err = domain.ValidateDescription(req.Description)
	if err != nil {
		ve.Add("description", err.Error())
	}
	if ve.HasErrors() {
		return domain.User{}, ve
	}

	user, err := s.repo.UpdateProfile(ctx, userID, req)
	if err != nil {
		return domain.User{}, fmt.Errorf("repo: %w", err)
	}

	user.AvatarUrl = s.avatarsBaseURL + "/" + user.AvatarUrl

	return user, nil
}

func (s *Service) SearchUsers(ctx context.Context, req domain.SearchUsersReq) ([]domain.User, error) {
	users, err := s.repo.SearchUsers(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("repo: %w", err)
	}

	for i := range users {
		users[i].AvatarUrl = s.avatarsBaseURL + "/" + users[i].AvatarUrl
	}

	return users, nil
}

func (s *Service) Follow(ctx context.Context, followerID uuid.UUID, followeeID uuid.UUID) error {
	err := s.repo.Follow(ctx, followerID, followeeID)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}
	return nil
}

func (s *Service) Unfollow(ctx context.Context, followerID uuid.UUID, followeeID uuid.UUID) error {
	err := s.repo.Unfollow(ctx, followerID, followeeID)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}
	return nil
}
