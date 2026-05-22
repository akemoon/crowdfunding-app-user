package domain

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/google/uuid"
)

type SearchUsersReq struct {
	Query  *string
	Limit  int
	Offset int
}

type UpdateProfileReq struct {
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
}

type User struct {
	ID             uuid.UUID `json:"id"`
	Username       string    `json:"username"`
	DisplayName    string    `json:"displayName"`
	Description    string    `json:"description"`
	AvatarUrl      string    `json:"avatarUrl"`
	FollowersCount int       `json:"followersCount"`
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
	if r == ' ' {
		return true
	}
	if unicode.Is(unicode.Latin, r) {
		return true
	}
	if unicode.Is(unicode.Cyrillic, r) {
		return true
	}
	if unicode.IsDigit(r) {
		return true
	}
	if unicode.IsPunct(r) {
		return true
	}
	return false
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
	if r == ' ' || r == '\n' {
		return true
	}
	if unicode.Is(unicode.Latin, r) {
		return true
	}
	if unicode.Is(unicode.Cyrillic, r) {
		return true
	}
	if unicode.IsDigit(r) {
		return true
	}
	if unicode.IsPunct(r) {
		return true
	}
	return false
}
