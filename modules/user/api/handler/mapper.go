package handler

import (
	"net/http"

	"github.com/akemoon/crowdfunding-app-user/modules/user/domain"
	"github.com/akemoon/golib/httplib"
)

var (
	ruleNotFound = httplib.ErrMapRule{
		Err:     domain.ErrNotFound,
		Status:  http.StatusNotFound,
		Code:    "not_found",
		Message: domain.ErrNotFound.Error(),
	}
	ruleUsernameExists = httplib.ErrMapRule{
		Err:     domain.ErrUsernameExists,
		Status:  http.StatusConflict,
		Code:    "username_exists",
		Message: domain.ErrUsernameExists.Error(),
	}
	ruleEmailExists = httplib.ErrMapRule{
		Err:     domain.ErrEmailExists,
		Status:  http.StatusConflict,
		Code:    "email_exists",
		Message: domain.ErrEmailExists.Error(),
	}
	ruleAlreadyFollowing = httplib.ErrMapRule{
		Err:     domain.ErrAlreadyFollowing,
		Status:  http.StatusConflict,
		Code:    "already_following",
		Message: domain.ErrAlreadyFollowing.Error(),
	}
	ruleFollowYourself = httplib.ErrMapRule{
		Err:     domain.ErrFollowToYourself,
		Status:  http.StatusUnprocessableEntity,
		Code:    "follow_yourself",
		Message: domain.ErrFollowToYourself.Error(),
	}
	ruleFollowerNotFound = httplib.ErrMapRule{
		Err:     domain.ErrFollowerNotFound,
		Status:  http.StatusNotFound,
		Code:    "follower_not_found",
		Message: domain.ErrFollowerNotFound.Error(),
	}
	ruleFolloweeNotFound = httplib.ErrMapRule{
		Err:     domain.ErrFolloweeNotFound,
		Status:  http.StatusNotFound,
		Code:    "followee_not_found",
		Message: domain.ErrFolloweeNotFound.Error(),
	}
)

var (
	createUserMapRules = []httplib.ErrMapRule{
		ruleUsernameExists,
		ruleEmailExists,
	}
	getUserMapRules = []httplib.ErrMapRule{
		ruleNotFound,
	}
	updateProfileMapRules = []httplib.ErrMapRule{
		ruleNotFound,
	}
	followMapRules = []httplib.ErrMapRule{
		ruleAlreadyFollowing,
		ruleFollowYourself,
		ruleFollowerNotFound,
		ruleFolloweeNotFound,
	}
)
