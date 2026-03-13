package domain

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

// TODO: rethink solution with lib

var (
	ErrInvalidUsername    = errors.New("invalid username")
	ErrUsernameExists     = errors.New("username already exists")
	ErrInvalidEmail       = errors.New("invalid email")
	ErrEmailExists        = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type CreateUserReq struct {
	Email        string
	Username     string
	PasswordHash string
}

type UserCredentials struct {
	UserID       uuid.UUID
	PasswordHash string
	Role         string
}

const (
	MinUsernameLength = 3
	MaxUsernameLength = 32
)

func ValidateUsername(username string) error {
	if len(username) < MinUsernameLength {
		return ErrInvalidUsername
	}
	if len(username) > MaxUsernameLength {
		return ErrInvalidUsername
	}
	for i := 0; i < len(username); i++ {
		if !isAllowedUsernameChar(username[i]) {
			return ErrInvalidUsername
		}
	}
	if username[0] == '-' || username[len(username)-1] == '-' {
		return ErrInvalidUsername
	}
	return nil
}

func isAllowedUsernameChar(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-'
}

// TODO: maybe move email to auth module
// TODO: tests
const MaxEmailLength = 254

func ValidateEmail(email string) error {
	if email == "" {
		return ErrInvalidEmail
	}
	if len(email) > MaxEmailLength {
		return ErrInvalidEmail
	}
	at := strings.LastIndex(email, "@")
	if at == -1 || at == 0 || at == len(email)-1 {
		return ErrInvalidEmail
	}
	local, domain := email[:at], email[at+1:]
	if len(local) > 64 {
		return ErrInvalidEmail
	}
	if !strings.Contains(domain, ".") {
		return ErrInvalidEmail
	}
	if strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") {
		return ErrInvalidEmail
	}
	return nil
}
