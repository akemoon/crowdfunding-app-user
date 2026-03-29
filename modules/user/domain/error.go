package domain

import "errors"

var (
	ErrUsernameExists = errors.New("username already exists")
	ErrEmailExists    = errors.New("email already exists")

	ErrInvalidDisplayName = errors.New("invalid display name")
	ErrInvalidDescription = errors.New("invalid description")

	ErrAlreadyFollowing = errors.New("already following")
	ErrFollowToYourself = errors.New("can't follow yourself")
	ErrFollowerNotFound = errors.New("follower not found")
	ErrFolloweeNotFound = errors.New("followee not found")

	ErrNotFound = errors.New("not found")
	ErrInternal = errors.New("internal error")
)
