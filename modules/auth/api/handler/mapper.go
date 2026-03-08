package handler

import (
	"net/http"

	"github.com/akemoon/crowdfunding-app-user/lib/api"
	"github.com/akemoon/crowdfunding-app-user/lib/httplib"
	"github.com/akemoon/crowdfunding-app-user/modules/auth/domain"
)

const (
	ErrCodeInvalidAccessToken  = "invalid_access_token"
	ErrCodeInvalidRefreshToken = "invalid_refresh_token"
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
)
