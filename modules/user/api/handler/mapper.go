package handler

import (
	"net/http"

	"github.com/akemoon/crowdfunding-app-user/lib/api"
	lib "github.com/akemoon/crowdfunding-app-user/lib/domain"
	"github.com/akemoon/golib/httplib"
	"github.com/akemoon/crowdfunding-app-user/modules/user/domain"
)

const (
	ErrCodeNotFound         = "not_found"
	ErrCodeAlreadyFollow    = "already_following"
	ErrCodeFollowYourself   = "follow_to_yourself"
	ErrCodeFollowerNotFound = "follower_not_found"
	ErrCodeFolloweeNotFound = "followee_not_found"
)

var (
	MapRuleErrNotFound = httplib.ErrMapRule{
		Err:     lib.ErrNotFound,
		Status:  http.StatusNotFound,
		Code:    ErrCodeNotFound,
		Message: lib.ErrNotFound.Error(),
	}
	MapRuleErrAlreadyFollowing = httplib.ErrMapRule{
		Err:     domain.ErrAlreadyFollowing,
		Status:  http.StatusConflict,
		Code:    ErrCodeAlreadyFollow,
		Message: domain.ErrAlreadyFollowing.Error(),
	}
	MapRuleErrFollowToYourself = httplib.ErrMapRule{
		Err:     domain.ErrFollowToYourself,
		Status:  http.StatusUnprocessableEntity,
		Code:    ErrCodeFollowYourself,
		Message: domain.ErrFollowToYourself.Error(),
	}
	MapRuleErrFollowerNotFound = httplib.ErrMapRule{
		Err:     domain.ErrFollowerNotFound,
		Status:  http.StatusNotFound,
		Code:    ErrCodeFollowerNotFound,
		Message: domain.ErrFollowerNotFound.Error(),
	}
	MapRuleErrFolloweeNotFound = httplib.ErrMapRule{
		Err:     domain.ErrFolloweeNotFound,
		Status:  http.StatusNotFound,
		Code:    ErrCodeFolloweeNotFound,
		Message: domain.ErrFolloweeNotFound.Error(),
	}
)

var (
	GetUserByIDMapRules = []httplib.ErrMapRule{
		MapRuleErrNotFound,
	}

	UpdateProfileMapRules = []httplib.ErrMapRule{
		api.MapRuleErrUsernameExists,
		MapRuleErrNotFound,
	}

	FollowMapRules = []httplib.ErrMapRule{
		MapRuleErrAlreadyFollowing,
		MapRuleErrFollowToYourself,
		MapRuleErrFollowerNotFound,
		MapRuleErrFolloweeNotFound,
	}
)
