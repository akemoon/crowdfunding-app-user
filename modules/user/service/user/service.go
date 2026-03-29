package user

import (
	"context"
	"fmt"

	lib "github.com/akemoon/crowdfunding-app-user/lib/domain"
	"github.com/akemoon/crowdfunding-app-user/modules/user/domain"
	userRepo "github.com/akemoon/crowdfunding-app-user/modules/user/repo/user"
	"github.com/akemoon/golib/validation"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo userRepo.Repo
}

func NewService(repo userRepo.Repo) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateUser(ctx context.Context, req domain.CreateUserReq) (uuid.UUID, error) {
	ve := &validation.Error{}

	if err := lib.ValidateUsername(req.Username); err != nil {
		ve.Add("username", err.Error())
	}
	if err := lib.ValidateEmail(req.Email); err != nil {
		ve.Add("email", err.Error())
	}
	if req.Password == "" {
		ve.Add("password", "password is required")
	}
	if ve.HasErrors() {
		return uuid.UUID{}, ve
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("hash password: %w", err)
	}

	id, err := s.repo.CreateUser(ctx, req, string(hash))
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("repo: %w", err)
	}

	return id, nil
}

func (s *Service) GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	u, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("repo: %w", err)
	}
	return u, nil
}

func (s *Service) ListUsers(ctx context.Context) ([]domain.User, error) {
	users, err := s.repo.ListUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("repo: %w", err)
	}

	if users == nil {
		users = []domain.User{}
	}

	return users, nil
}

func (s *Service) UpdateProfile(ctx context.Context, userID uuid.UUID, req domain.UpdateProfileReq) (domain.User, error) {
	ve := &validation.Error{}

	if err := domain.ValidateDisplayName(req.DisplayName); err != nil {
		ve.Add("displayName", err.Error())
	}
	if err := domain.ValidateDescription(req.Description); err != nil {
		ve.Add("description", err.Error())
	}
	if ve.HasErrors() {
		return domain.User{}, ve
	}

	u, err := s.repo.UpdateProfile(ctx, userID, req)
	if err != nil {
		return domain.User{}, fmt.Errorf("repo: %w", err)
	}
	return u, nil
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

func (s *Service) GetFollows(ctx context.Context, userID uuid.UUID) ([]domain.User, error) {
	users, err := s.repo.GetFollows(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("repo: %w", err)
	}
	if users == nil {
		users = []domain.User{}
	}
	return users, nil
}
