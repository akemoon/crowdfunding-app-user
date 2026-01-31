package domain

import "errors"

var (
	ErrInvalidUsername = errors.New("invalid username")
	ErrUsernameExists  = errors.New("username already exists")
	ErrUnknownConflict = errors.New("unknown conflict")
	ErrUserNotFound    = errors.New("user not found")

	ErrInternal = errors.New("internal error")
)
