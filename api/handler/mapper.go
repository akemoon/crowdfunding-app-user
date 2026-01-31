package handler

import (
	"errors"
	"net/http"

	"github.com/akemoon/crowdfunding-app-user/domain"
)

const (
	HttpErrInvalidUsername = "invalid_username"
	HttpErrUsernameExists  = "username_exists"
	HttpErrUnknownConflict = "unknown_conflict"
	HttpErrUserNotFound    = "user_not_found"

	HttpErrInternal = "internal_error"
)

type ErrResp struct {
	Error   string `json:"error"`
	Details string `json:"details"`
}

func MapErrToHTTP(err error) (int, ErrResp) {
	if errors.Is(err, domain.ErrInvalidUsername) {
		return http.StatusBadRequest, ErrResp{
			Error:   HttpErrInvalidUsername,
			Details: domain.ErrInvalidUsername.Error(),
		}
	}

	if errors.Is(err, domain.ErrUsernameExists) {
		return http.StatusConflict, ErrResp{
			Error:   HttpErrUsernameExists,
			Details: domain.ErrUsernameExists.Error(),
		}
	}

	if errors.Is(err, domain.ErrUnknownConflict) {
		return http.StatusConflict, ErrResp{
			Error:   HttpErrUnknownConflict,
			Details: domain.ErrUnknownConflict.Error(),
		}
	}

	if errors.Is(err, domain.ErrUserNotFound) {
		return http.StatusNotFound, ErrResp{
			Error:   HttpErrUserNotFound,
			Details: domain.ErrUserNotFound.Error(),
		}
	}

	return http.StatusInternalServerError, ErrResp{
		Error:   HttpErrInternal,
		Details: domain.ErrInternal.Error(),
	}
}
