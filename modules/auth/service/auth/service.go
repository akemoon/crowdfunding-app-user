package auth

import (
	"context"
	"fmt"
	"log"

	lib "github.com/akemoon/crowdfunding-app-user/lib/domain"
	"github.com/akemoon/crowdfunding-app-user/modules/auth/domain"
	"github.com/akemoon/crowdfunding-app-user/modules/auth/repo/user"
	"github.com/akemoon/crowdfunding-app-user/modules/auth/service/token"
	"github.com/akemoon/crowdfunding-app-user/modules/auth/tool/hasher"
	userPublisher "github.com/akemoon/crowdfunding-app-user/platform/publisher/user"
	"github.com/akemoon/golib/validation"
	"github.com/google/uuid"
)

type Service struct {
	userRepo  user.Repo
	hasher    hasher.Hasher
	tokenSvc  *token.Service // TODO: interface
	publisher *userPublisher.Publisher
}

func NewService(r user.Repo, h hasher.Hasher, t *token.Service, p *userPublisher.Publisher) *Service {
	return &Service{
		userRepo:  r,
		hasher:    h,
		tokenSvc:  t,
		publisher: p,
	}
}

func (s *Service) SignUp(ctx context.Context, req domain.SignUpReq) error {
	ve := &validation.Error{}

	if err := lib.ValidateEmail(req.Email); err != nil {
		ve.Add("email", err.Error())
	}
	if err := lib.ValidateUsername(req.Username); err != nil {
		ve.Add("username", err.Error())
	}
	if err := domain.ValidatePassword(req.Password); err != nil {
		ve.Add("password", err.Error())
	}
	if ve.HasErrors() {
		return ve
	}

	passwordHash, err := s.hasher.Hash(req.Password)
	if err != nil {
		return err
	}

	userID, err := s.userRepo.CreateUser(ctx, lib.CreateUserReq{
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: passwordHash,
	})
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}

	// TODO: outbox pattern to guarantee delivery
	event := userPublisher.Event{
		EventID: uuid.New(),
		Type:    userPublisher.EventTypeRegistered,
		UserID:  userID,
	}
	if err := s.publisher.Publish(ctx, event); err != nil {
		log.Printf("auth: publish user registered event: %v", err)
	}

	return nil
}

func (s *Service) SignIn(ctx context.Context, req domain.SignInReq) (domain.SignInResp, error) {
	creds, err := s.userRepo.GetCredentialsByEmail(ctx, req.Email)
	if err != nil {
		return domain.SignInResp{}, err
	}

	err = s.hasher.Compare(req.Password, creds.PasswordHash)
	if err != nil {
		return domain.SignInResp{}, lib.ErrInvalidCredentials
	}

	claims := domain.TokenClaims{UserID: creds.UserID, Role: creds.Role}

	accessToken, err := s.tokenSvc.GenerateAccessToken(claims)
	if err != nil {
		return domain.SignInResp{}, fmt.Errorf("token service: %w", err)
	}

	refreshToken, err := s.tokenSvc.GenerateRefreshToken(ctx, creds.UserID)
	if err != nil {
		return domain.SignInResp{}, fmt.Errorf("token service: %w", err)
	}

	return domain.SignInResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) SignOut(ctx context.Context, req domain.SignOutReq) error {
	if req.RefreshToken == "" {
		return domain.ErrInvalidRefreshToken
	}

	err := s.tokenSvc.DeleteRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return fmt.Errorf("token service: %w", err)
	}

	return nil
}

func (s *Service) ValidateAccessToken(token string) (domain.TokenClaims, error) {
	claims, err := s.tokenSvc.ValidateAccessToken(token)
	if err != nil {
		return domain.TokenClaims{}, err
	}

	return claims, nil
}

func (s *Service) Refresh(ctx context.Context, req domain.RefreshReq) (domain.RefreshResp, error) {
	if req.RefreshToken == "" {
		return domain.RefreshResp{}, domain.ErrInvalidRefreshToken
	}

	userID, err := s.tokenSvc.ValidateRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return domain.RefreshResp{}, err
	}

	creds, err := s.userRepo.GetCredentialsByID(ctx, userID)
	if err != nil {
		return domain.RefreshResp{}, fmt.Errorf("repo: %w", err)
	}

	err = s.tokenSvc.DeleteRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return domain.RefreshResp{}, fmt.Errorf("token service: %w", err)
	}

	newRefreshToken, err := s.tokenSvc.GenerateRefreshToken(ctx, userID)
	if err != nil {
		return domain.RefreshResp{}, fmt.Errorf("token service: %w", err)
	}

	claims := domain.TokenClaims{UserID: userID, Role: creds.Role}

	newAccessToken, err := s.tokenSvc.GenerateAccessToken(claims)
	if err != nil {
		return domain.RefreshResp{}, fmt.Errorf("token service: %w", err)
	}

	return domain.RefreshResp{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

// TODO: block rule: delete refresh tokens from db
