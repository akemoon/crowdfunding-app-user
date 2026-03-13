package token

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/akemoon/crowdfunding-app-user/modules/auth/domain"
	"github.com/akemoon/crowdfunding-app-user/modules/auth/repo/token"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	accessTokenLifeTime  = 15 * time.Minute
	refreshTokenLifeTime = 7 * 24 * time.Hour
)

type Service struct {
	refreshTokenRepo token.Repo
	secret           string
}

func NewService(r token.Repo, s string) *Service {
	return &Service{
		refreshTokenRepo: r,
		secret:           s,
	}
}

func (s *Service) GenerateAccessToken(tc domain.TokenClaims) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": tc.UserID,
		"role":   tc.Role,
		"exp":    time.Now().Add(accessTokenLifeTime).Unix(),
	})

	signed, err := t.SignedString([]byte(s.secret))
	if err != nil {
		return "", fmt.Errorf("sign access token: %w", err)
	}

	return signed, nil
}

func (s *Service) GenerateRefreshToken(ctx context.Context, tc domain.TokenClaims) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": tc.UserID,
		"exp":    time.Now().Add(refreshTokenLifeTime).Unix(),
	})

	signed, err := t.SignedString([]byte(s.secret))
	if err != nil {
		return "", fmt.Errorf("sign refresh token: %w", err)
	}

	err = s.refreshTokenRepo.Set(ctx, signed, refreshTokenLifeTime)
	if err != nil {
		return "", fmt.Errorf("token repo: %w", err)
	}

	return signed, nil
}

func (s *Service) DeleteRefreshToken(ctx context.Context, refreshToken string) error {
	err := s.refreshTokenRepo.Delete(ctx, refreshToken)
	if err != nil {
		return fmt.Errorf("token repo: %w", err)
	}

	return nil
}

func (s *Service) ValidateAccessToken(token string) (domain.TokenClaims, error) {
	parts := strings.Fields(token)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return domain.TokenClaims{}, domain.ErrInvalidAccessToken
	}

	parser := jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	jwtToken, err := parser.Parse(parts[1], func(t *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})
	if err != nil {
		return domain.TokenClaims{}, domain.ErrInvalidAccessToken
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return domain.TokenClaims{}, domain.ErrInvalidAccessToken
	}

	userIDStr, ok := claims["userID"].(string)
	if !ok || userIDStr == "" {
		return domain.TokenClaims{}, domain.ErrInvalidAccessToken
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return domain.TokenClaims{}, domain.ErrInvalidAccessToken
	}

	role, ok := claims["role"].(string)
	if !ok || role == "" {
		return domain.TokenClaims{}, domain.ErrInvalidAccessToken
	}

	return domain.TokenClaims{UserID: userID, Role: role}, nil
}
