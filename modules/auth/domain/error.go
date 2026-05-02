package domain

import "errors"

var (
	ErrInvalidPassword = errors.New("invalid password")

	ErrUnknownConflict = errors.New("unknown conflict")

	ErrInvalidAccessToken = errors.New("invalid access token")

	ErrInvalidRefreshToken = errors.New("invalid refresh token")

	ErrInternal = errors.New("internal error")

	ErrInvalidRole = errors.New("invalid role")

	ErrForbidden = errors.New("forbidden")
)
