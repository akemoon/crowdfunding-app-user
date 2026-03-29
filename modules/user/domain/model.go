package domain

import (
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/google/uuid"
)

// CreateUserReq — тело запроса POST /users
type CreateUserReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UpdateProfileReq — тело запроса PUT /users/:id
type UpdateProfileReq struct {
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
}

// FollowReq — тело запроса POST/DELETE /users/:id/follow
type FollowReq struct {
	FollowerID uuid.UUID `json:"followerId"`
}

// User — модель пользователя в ответе
type User struct {
	ID          uuid.UUID  `json:"id"`
	Username    string     `json:"username"`
	DisplayName string     `json:"displayName"`
	Description string     `json:"description"`
	AvatarKey   string     `json:"avatarKey"`
	UpdatedAt   *time.Time `json:"updatedAt"`
}

const MaxDisplayNameLength = 120

func ValidateDisplayName(displayName string) error {
	if displayName == "" {
		return nil
	}
	if utf8.RuneCountInString(displayName) > MaxDisplayNameLength {
		return ErrInvalidDisplayName
	}
	if displayName != strings.TrimSpace(displayName) {
		return ErrInvalidDisplayName
	}
	for _, r := range displayName {
		if !isAllowedDisplayNameRune(r) {
			return ErrInvalidDisplayName
		}
	}
	return nil
}

func isAllowedDisplayNameRune(r rune) bool {
	return r == ' ' ||
		unicode.Is(unicode.Latin, r) ||
		unicode.Is(unicode.Cyrillic, r) ||
		unicode.IsDigit(r) ||
		unicode.IsPunct(r)
}

const MaxDescriptionLength = 200

func ValidateDescription(description string) error {
	if description == "" {
		return nil
	}
	if description != strings.TrimSpace(description) {
		return ErrInvalidDescription
	}
	if utf8.RuneCountInString(description) > MaxDescriptionLength {
		return ErrInvalidDescription
	}
	for _, r := range description {
		if !isAllowedDescriptionRune(r) {
			return ErrInvalidDescription
		}
	}
	return nil
}

func isAllowedDescriptionRune(r rune) bool {
	return r == ' ' || r == '\n' ||
		unicode.Is(unicode.Latin, r) ||
		unicode.Is(unicode.Cyrillic, r) ||
		unicode.IsDigit(r) ||
		unicode.IsPunct(r)
}
