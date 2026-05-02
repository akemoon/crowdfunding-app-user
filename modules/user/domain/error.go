package domain

import "errors"

var (
	ErrInvalidDisplayName = errors.New("invalid display name")
	ErrInvalidDescription = errors.New("invalid description")

	ErrAlreadyFollowing = errors.New("already following")
	ErrFollowToYourself = errors.New("can't follow to yourself")
	ErrFollowerNotFound = errors.New("follower not found")
	ErrFolloweeNotFound = errors.New("followee not found")

	ErrUnknownConflict = errors.New("unknown conflict")

	ErrInternal = errors.New("internal error")
)
