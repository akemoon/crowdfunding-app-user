package handler

import (
	"net/http"

	"github.com/akemoon/crowdfunding-app-user/lib/api"
	lib "github.com/akemoon/crowdfunding-app-user/lib/domain"
	"github.com/akemoon/golib/httplib"
	"github.com/akemoon/crowdfunding-app-user/modules/auth/domain"
)

const (
	ErrCodeInvalidAccessToken  = "invalid_access_token"
	ErrCodeInvalidRefreshToken = "invalid_refresh_token"
	ErrCodeInvalidRole         = "invalid_role"
	ErrCodeForbidden           = "forbidden"
	ErrCodeNotFound            = "not_found"
)

var (
	SignUpMapRules = []httplib.ErrMapRule{
		api.MapRuleErrUsernameExists,
		api.MapRuleErrEmailExists,
	}

	SignInMapRules = []httplib.ErrMapRule{
		api.MapRuleErrInvalidCredentials,
	}

	SignOutMapRules = []httplib.ErrMapRule{
		{
			Err:     domain.ErrInvalidRefreshToken,
			Status:  http.StatusBadRequest,
			Code:    ErrCodeInvalidRefreshToken,
			Message: domain.ErrInvalidRefreshToken.Error(),
		},
	}

	CheckMapRules = []httplib.ErrMapRule{
		{
			Err:     domain.ErrInvalidAccessToken,
			Status:  http.StatusUnauthorized,
			Code:    ErrCodeInvalidAccessToken,
			Message: domain.ErrInvalidAccessToken.Error(),
		},
	}

	RefreshMapRules = []httplib.ErrMapRule{
		{
			Err:     domain.ErrInvalidRefreshToken,
			Status:  http.StatusUnauthorized,
			Code:    ErrCodeInvalidRefreshToken,
			Message: domain.ErrInvalidRefreshToken.Error(),
		},
		{
			Err:     lib.ErrNotFound,
			Status:  http.StatusUnauthorized,
			Code:    ErrCodeInvalidRefreshToken,
			Message: domain.ErrInvalidRefreshToken.Error(),
		},
	}

	GetCredentialsByIDMapRules = []httplib.ErrMapRule{
		{
			Err:     domain.ErrForbidden,
			Status:  http.StatusForbidden,
			Code:    ErrCodeForbidden,
			Message: domain.ErrForbidden.Error(),
		},
		{
			Err:     lib.ErrNotFound,
			Status:  http.StatusNotFound,
			Code:    ErrCodeNotFound,
			Message: lib.ErrNotFound.Error(),
		},
	}

	UpdateRoleMapRules = []httplib.ErrMapRule{
		{
			Err:     domain.ErrForbidden,
			Status:  http.StatusForbidden,
			Code:    ErrCodeForbidden,
			Message: domain.ErrForbidden.Error(),
		},
		{
			Err:     domain.ErrInvalidRole,
			Status:  http.StatusBadRequest,
			Code:    ErrCodeInvalidRole,
			Message: domain.ErrInvalidRole.Error(),
		},
	}
)
