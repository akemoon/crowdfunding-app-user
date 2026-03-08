package auth

import (
	"context"
	"fmt"

	lib "github.com/akemoon/crowdfunding-app-user/lib/domain"
	"github.com/akemoon/crowdfunding-app-user/lib/validation"
	"github.com/akemoon/crowdfunding-app-user/modules/auth/domain"
	"github.com/akemoon/crowdfunding-app-user/modules/auth/repo/user"
	"github.com/akemoon/crowdfunding-app-user/modules/auth/service/token"
	"github.com/akemoon/crowdfunding-app-user/modules/auth/tool/hasher"
)

type Service struct {
	userRepo user.Repo
	hasher   hasher.Hasher
	tokenSvc *token.Service // TODO: interface
}

func NewService(r user.Repo, h hasher.Hasher, t *token.Service) *Service {
	return &Service{
		userRepo: r,
		hasher:   h,
		tokenSvc: t,
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

	err = s.userRepo.CreateUser(ctx, lib.CreateUserReq{
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: passwordHash,
	})
	if err != nil {
		return fmt.Errorf("repo: %w", err)
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

	claims := domain.TokenClaims{UserID: creds.UserID}

	accessToken, err := s.tokenSvc.GenerateAccessToken(claims)
	if err != nil {
		return domain.SignInResp{}, fmt.Errorf("token service: %w", err)
	}

	refreshToken, err := s.tokenSvc.GenerateRefreshToken(ctx, claims)
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
	userID, err := s.tokenSvc.ValidateAccessToken(token)
	if err != nil {
		return domain.TokenClaims{}, err
	}

	return domain.TokenClaims{UserID: userID}, nil
}

// TODO: block rule: delete refresh tokens from db
