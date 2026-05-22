package domain

import "github.com/google/uuid"

type SignUpReq struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignInReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInResp struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type SignOutReq struct {
	RefreshToken string `json:"refreshToken"`
}

type RefreshReq struct {
	RefreshToken string `json:"refreshToken"`
}

type RefreshResp struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

const (
	RoleUser  = "user"
	RoleModer = "moder"
	RoleAdmin = "admin"
)

type CredentialsResp struct {
	UserID    uuid.UUID `json:"userID"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	IsBlocked bool      `json:"isBlocked"`
}

type UpdateRoleReq struct {
	UserID     uuid.UUID
	NewRole    string
	CallerRole string
}

func ValidateRole(role string) error {
	switch role {
	case RoleUser, RoleModer, RoleAdmin:
		return nil
	default:
		return ErrInvalidRole
	}
}

const (
	MinPasswordLength = 12
	MaxPasswordLength = 64
)

func ValidatePassword(password string) error {
	if len(password) < MinPasswordLength {
		return ErrInvalidPassword
	}
	if len(password) > MaxPasswordLength {
		return ErrInvalidPassword
	}

	// NOTE: NIST recommends length over complexity, but classic rules also kept
	// https://pages.nist.gov/800-63-4/sp800-63b.html
	//
	var hasLower, hasUpper, hasDigit bool

	for i := 0; i < len(password); i++ {
		c := password[i]

		// Only ASCII
		if c < ' ' || c > '~' {
			return ErrInvalidPassword
		}

		switch {
		case c >= 'a' && c <= 'z':
			hasLower = true
		case c >= 'A' && c <= 'Z':
			hasUpper = true
		case c >= '0' && c <= '9':
			hasDigit = true
		}
	}

	if !hasLower || !hasUpper || !hasDigit {
		return ErrInvalidPassword
	}

	return nil
}
