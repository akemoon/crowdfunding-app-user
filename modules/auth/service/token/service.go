package token

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/akemoon/crowdfunding-app-user/modules/auth/domain"
	"github.com/akemoon/crowdfunding-app-user/modules/auth/repo/refresh"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	accessTokenLifeTime  = 15 * time.Minute
	refreshTokenLifeTime = 24 * 60 * time.Minute
)

type Service struct {
	refreshRepo refresh.Repo
	secret      string
}

func NewService(r refresh.Repo, s string) *Service {
	return &Service{
		refreshRepo: r,
		secret:      s,
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

func (s *Service) GenerateRefreshToken(ctx context.Context, userID uuid.UUID) (string, error) {
	raw := make([]byte, 32)

	_, err := rand.Read(raw)
	if err != nil {
		return "", fmt.Errorf("generate refresh token: %w", err)
	}

	tok := hex.EncodeToString(raw)
	hash := sha256hex(tok)

	err = s.refreshRepo.Set(ctx, hash, userID, refreshTokenLifeTime)
	if err != nil {
		return "", fmt.Errorf("token repo: %w", err)
	}

	return tok, nil
}

func (s *Service) ValidateRefreshToken(ctx context.Context, tok string) (uuid.UUID, error) {
	hash := sha256hex(tok)

	userID, err := s.refreshRepo.Get(ctx, hash)
	if err != nil {
		return uuid.UUID{}, domain.ErrInvalidRefreshToken
	}

	return userID, nil
}

func (s *Service) DeleteRefreshToken(ctx context.Context, tok string) error {
	hash := sha256hex(tok)

	err := s.refreshRepo.Delete(ctx, hash)
	if err != nil {
		return fmt.Errorf("token repo: %w", err)
	}

	return nil
}

func (s *Service) DeleteAllRefreshTokensByUserID(ctx context.Context, userID uuid.UUID) error {
	err := s.refreshRepo.DeleteAllByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("token repo: %w", err)
	}

	return nil
}

func (s *Service) ValidateAccessToken(tok string) (domain.TokenClaims, error) {
	parts := strings.Fields(tok)
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

func sha256hex(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}
