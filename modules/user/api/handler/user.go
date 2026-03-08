package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/akemoon/crowdfunding-app-user/lib/httplib"
	"github.com/akemoon/crowdfunding-app-user/modules/user/domain"
	"github.com/akemoon/crowdfunding-app-user/modules/user/service/user"
	"github.com/google/uuid"
)

const userIDHeader = "X-User-ID"

// @Summary Get user profile
// @Description Get user profile by path user id
// @Accept json
// @Produce json
// @Param id path string true "User id (uuid)"
// @Success 200 {object} domain.User "User profile"
// @Failure 400 {string} string "Invalid path id"
// @Failure 404 {object} httplib.ErrResp "User not found"
// @Failure 405 {string} string "Method not allowed"
// @Failure 500 {object} httplib.ErrResp "Internal server error"
// @Router /users/{id} [get]
func GetUserByID(svc *user.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			http.Error(w, "invalid path id", http.StatusBadRequest)
			return
		}

		user, err := svc.GetUserByID(r.Context(), id)
		if err != nil {
			log.Println(err)
			status, errResp := httplib.MapErrToHTTP(err, GetUserByIDMapRules)
			httplib.WriteJSON(w, status, errResp)
			return
		}

		httplib.WriteJSON(w, http.StatusOK, user)
	}
}

// @Summary Update user profile
// @Description Update user profile by given payload
// @Accept json
// @Produce json
// @Param X-User-ID header string true "User id (uuid)"
// @Param payload body domain.UpdateProfileReq true "Update user profile payload"
// @Success 200 {object} domain.User "User profile"
// @Failure 400 {string} string "Invalid request body"
// @Failure 400 {object} httplib.ErrResp "Validation error"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {object} httplib.ErrResp "User not found"
// @Failure 405 {string} string "Method not allowed"
// @Failure 409 {object} httplib.ErrResp "Conflict"
// @Failure 500 {object} httplib.ErrResp "Internal server error"
// @Router /users/me/profile [patch]
func UpdateProfile(svc *user.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := httplib.ParseUUIDHeader(r, userIDHeader)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		var req domain.UpdateProfileReq

		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		user, err := svc.UpdateProfile(r.Context(), id, req)
		if err != nil {
			log.Println(err)
			status, errResp := httplib.MapErrToHTTP(err, UpdateProfileMapRules)
			httplib.WriteJSON(w, status, errResp)
			return
		}

		httplib.WriteJSON(w, http.StatusOK, user)
	}
}

// @Summary Follow to a user
// @Description Follow to a user by id
// @Accept json
// @Produce json
// @Param X-User-ID header string true "Follower id (uuid)"
// @Param id path string true "Followee id (uuid)"
// @Success 201 {string} string "Created"
// @Failure 400 {string} string "Invalid path id"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {object} httplib.ErrResp "Follower or followee not found"
// @Failure 405 {string} string "Method not allowed"
// @Failure 409 {object} httplib.ErrResp "Already following"
// @Failure 422 {object} httplib.ErrResp "Cannot follow yourself"
// @Failure 500 {object} httplib.ErrResp "Internal server error"
// @Router /users/{id}/follow [post]
func Follow(svc *user.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		followerID, err := httplib.ParseUUIDHeader(r, userIDHeader)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		followeeID, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			http.Error(w, "invalid path id", http.StatusBadRequest)
			return
		}

		err = svc.Follow(r.Context(), followerID, followeeID)
		if err != nil {
			log.Println(err)
			status, errResp := httplib.MapErrToHTTP(err, FollowMapRules)
			httplib.WriteJSON(w, status, errResp)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

// @Summary Unfollow to a user
// @Description Unfollow to a user by id
// @Accept json
// @Produce json
// @Param X-User-ID header string true "Follower id (uuid)"
// @Param id path string true "Followee id (uuid)"
// @Success 204 {string} string "No content"
// @Failure 400 {string} string "Invalid path id"
// @Failure 401 {string} string "Unauthorized"
// @Failure 405 {string} string "Method not allowed"
// @Failure 500 {object} string "Internal server error"
// @Router /users/{id}/follow [delete]
func Unfollow(svc *user.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		followerID, err := httplib.ParseUUIDHeader(r, userIDHeader)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		followeeID, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			http.Error(w, "invalid path id", http.StatusBadRequest)
			return
		}

		err = svc.Unfollow(r.Context(), followerID, followeeID)
		if err != nil {
			log.Println(err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// TODO: implement
//
// @Router /me/subscriptions [get]
// @Router /user/{id}/subscriptions [get]
