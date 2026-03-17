package api

import (
	"net/http"

	lib "github.com/akemoon/crowdfunding-app-user/lib/domain"
	"github.com/akemoon/golib/httplib"
)

const (
	ErrCodeUsernameExists     = "username_exists"
	ErrCodeEmailExists        = "email_exists"
	ErrCodeInvalidCredentials = "invalid_credentials"
)

var (
	MapRuleErrUsernameExists = httplib.ErrMapRule{
		Err:     lib.ErrUsernameExists,
		Status:  http.StatusConflict,
		Code:    ErrCodeUsernameExists,
		Message: lib.ErrUsernameExists.Error(),
	}

	MapRuleErrEmailExists = httplib.ErrMapRule{
		Err:     lib.ErrEmailExists,
		Status:  http.StatusConflict,
		Code:    ErrCodeEmailExists,
		Message: lib.ErrEmailExists.Error(),
	}

	MapRuleErrInvalidCredentials = httplib.ErrMapRule{
		Err:     lib.ErrInvalidCredentials,
		Status:  http.StatusUnauthorized,
		Code:    ErrCodeInvalidCredentials,
		Message: lib.ErrInvalidCredentials.Error(),
	}
)
