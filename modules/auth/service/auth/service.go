package auth

import (
	"context"
	"fmt"
	"log"

	lib "github.com/akemoon/crowdfunding-app-user/lib/domain"
	"github.com/akemoon/crowdfunding-app-user/modules/auth/domain"
	authrepo "github.com/akemoon/crowdfunding-app-user/modules/auth/repo/auth"
	"github.com/akemoon/crowdfunding-app-user/modules/auth/service/token"
	"github.com/akemoon/crowdfunding-app-user/modules/auth/tool/hasher"
	userPublisher "github.com/akemoon/crowdfunding-app-user/platform/publisher/user"
	"github.com/akemoon/golib/validation"
	"github.com/google/uuid"
)

type Service struct {
	userRepo  authrepo.Repo
	hasher    hasher.Hasher
	tokenSvc  *token.Service // TODO: interface
	publisher *userPublisher.Publisher
}

func NewService(r authrepo.Repo, h hasher.Hasher, t *token.Service, p *userPublisher.Publisher) *Service {
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

	if creds.IsBlocked {
		return domain.SignInResp{}, domain.ErrUserBlocked
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

func (s *Service) GetMe(ctx context.Context, userID uuid.UUID) (domain.GetMeResp, error) {
	creds, err := s.userRepo.GetCredentialsByID(ctx, userID)
	if err != nil {
		return domain.GetMeResp{}, fmt.Errorf("repo: %w", err)
	}

	return domain.GetMeResp{
		Email: creds.Email,
		Role:  creds.Role,
	}, nil
}

func (s *Service) CheckAccess(authHeader string) (domain.TokenClaims, error) {
	claims, err := s.tokenSvc.ValidateAccessToken(authHeader)
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

func (s *Service) GetCredentialsByID(ctx context.Context, callerRole string, userID uuid.UUID) (domain.CredentialsResp, error) {
	if callerRole != domain.RoleAdmin {
		return domain.CredentialsResp{}, domain.ErrForbidden
	}

	creds, err := s.userRepo.GetCredentialsByID(ctx, userID)
	if err != nil {
		return domain.CredentialsResp{}, fmt.Errorf("repo: %w", err)
	}

	return domain.CredentialsResp{
		UserID:    creds.UserID,
		Email:     creds.Email,
		Role:      creds.Role,
		IsBlocked: creds.IsBlocked,
	}, nil
}

func (s *Service) UpdateRole(ctx context.Context, req domain.UpdateRoleReq) error {
	if req.CallerRole != domain.RoleAdmin {
		return domain.ErrForbidden
	}

	err := domain.ValidateRole(req.NewRole)
	if err != nil {
		return err
	}

	err = s.userRepo.UpdateRole(ctx, req.UserID, req.NewRole)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}

	return nil
}

func (s *Service) SetBlocked(ctx context.Context, callerRole string, userID uuid.UUID, blocked bool) error {
	if callerRole != domain.RoleAdmin {
		return domain.ErrForbidden
	}

	err := s.userRepo.SetBlocked(ctx, userID, blocked)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}

	if blocked {
		err = s.tokenSvc.DeleteAllRefreshTokensByUserID(ctx, userID)
		if err != nil {
			return fmt.Errorf("token service: %w", err)
		}
	}

	return nil
}
