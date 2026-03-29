package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/akemoon/crowdfunding-app-user/modules/user/domain"
	"github.com/akemoon/crowdfunding-app-user/modules/user/service/user"
	"github.com/akemoon/golib/httplib"
	"github.com/google/uuid"
)

// ErrorResponse is a standard error response.
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// CreateUser godoc
//
//	@Summary      Create user
//	@Description  Creates a new user and their credentials (users + credentials)
//	@Tags         users
//	@Accept       json
//	@Produce      json
//	@Param        body  body      domain.CreateUserReq  true  "New user data"
//	@Success      201   {object}  domain.User
//	@Failure      400   {object}  ErrorResponse
//	@Failure      409   {object}  ErrorResponse  "username or email already taken"
//	@Failure      500   {object}  ErrorResponse
//	@Router       /users [post]
func CreateUser(svc *user.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req domain.CreateUserReq

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		id, err := svc.CreateUser(r.Context(), req)
		if err != nil {
			log.Println(err)
			status, errResp := httplib.MapErrToHTTP(err, createUserMapRules)
			httplib.WriteJSON(w, status, errResp)
			return
		}

		u, err := svc.GetUserByID(r.Context(), id)
		if err != nil {
			log.Println(err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		httplib.WriteJSON(w, http.StatusCreated, u)
	}
}

// GetUserByID godoc
//
//	@Summary      Get user by ID
//	@Description  Returns user profile by UUID
//	@Tags         users
//	@Produce      json
//	@Param        id   path      string  true  "User UUID"
//	@Success      200  {object}  domain.User
//	@Failure      400  {object}  ErrorResponse
//	@Failure      404  {object}  ErrorResponse
//	@Failure      500  {object}  ErrorResponse
//	@Router       /users/{id} [get]
func GetUserByID(svc *user.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			http.Error(w, "invalid path id", http.StatusBadRequest)
			return
		}

		u, err := svc.GetUserByID(r.Context(), id)
		if err != nil {
			log.Println(err)
			status, errResp := httplib.MapErrToHTTP(err, getUserMapRules)
			httplib.WriteJSON(w, status, errResp)
			return
		}

		httplib.WriteJSON(w, http.StatusOK, u)
	}
}

// ListUsers godoc
//
//	@Summary      List users
//	@Description  Returns all users
//	@Tags         users
//	@Produce      json
//	@Success      200  {array}   domain.User
//	@Failure      500  {object}  ErrorResponse
//	@Router       /users [get]
func ListUsers(svc *user.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := svc.ListUsers(r.Context())
		if err != nil {
			log.Println(err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		httplib.WriteJSON(w, http.StatusOK, users)
	}
}

// UpdateProfile godoc
//
//	@Summary      Update profile
//	@Description  Updates display_name and description of the user
//	@Tags         users
//	@Accept       json
//	@Produce      json
//	@Param        id    path      string                   true  "User UUID"
//	@Param        body  body      domain.UpdateProfileReq  true  "New profile data"
//	@Success      200   {object}  domain.User
//	@Failure      400   {object}  ErrorResponse
//	@Failure      404   {object}  ErrorResponse
//	@Failure      500   {object}  ErrorResponse
//	@Router       /users/{id} [put]
func UpdateProfile(svc *user.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			http.Error(w, "invalid path id", http.StatusBadRequest)
			return
		}

		var req domain.UpdateProfileReq

		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		u, err := svc.UpdateProfile(r.Context(), id, req)
		if err != nil {
			log.Println(err)
			status, errResp := httplib.MapErrToHTTP(err, updateProfileMapRules)
			httplib.WriteJSON(w, status, errResp)
			return
		}

		httplib.WriteJSON(w, http.StatusOK, u)
	}
}

// Follow godoc
//
//	@Summary      Follow user
//	@Description  Creates a record in the follows table (follower -> followee)
//	@Tags         follows
//	@Accept       json
//	@Param        id    path  string            true  "Followee UUID"
//	@Param        body  body  domain.FollowReq  true  "Follower UUID"
//	@Success      201
//	@Failure      400  {object}  ErrorResponse
//	@Failure      404  {object}  ErrorResponse
//	@Failure      409  {object}  ErrorResponse  "already following"
//	@Failure      422  {object}  ErrorResponse  "cannot follow yourself"
//	@Failure      500  {object}  ErrorResponse
//	@Router       /users/{id}/follow [post]
func Follow(svc *user.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		followeeID, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			http.Error(w, "invalid path id", http.StatusBadRequest)
			return
		}

		var req domain.FollowReq

		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		err = svc.Follow(r.Context(), req.FollowerID, followeeID)
		if err != nil {
			log.Println(err)
			status, errResp := httplib.MapErrToHTTP(err, followMapRules)
			httplib.WriteJSON(w, status, errResp)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

// GetFollows godoc
//
//	@Summary      Get follows
//	@Description  Returns the list of users that the given user follows
//	@Tags         follows
//	@Produce      json
//	@Param        id   path      string  true  "User UUID"
//	@Success      200  {array}   domain.User
//	@Failure      400  {object}  ErrorResponse
//	@Failure      500  {object}  ErrorResponse
//	@Router       /users/{id}/follows [get]
func GetFollows(svc *user.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			http.Error(w, "invalid path id", http.StatusBadRequest)
			return
		}

		users, err := svc.GetFollows(r.Context(), id)
		if err != nil {
			log.Println(err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		httplib.WriteJSON(w, http.StatusOK, users)
	}
}

// Unfollow godoc
//
//	@Summary      Unfollow user
//	@Description  Removes a record from the follows table
//	@Tags         follows
//	@Accept       json
//	@Param        id    path  string            true  "Followee UUID"
//	@Param        body  body  domain.FollowReq  true  "Follower UUID"
//	@Success      204
//	@Failure      400  {object}  ErrorResponse
//	@Failure      500  {object}  ErrorResponse
//	@Router       /users/{id}/follow [delete]
func Unfollow(svc *user.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		followeeID, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			http.Error(w, "invalid path id", http.StatusBadRequest)
			return
		}

		var req domain.FollowReq

		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		err = svc.Unfollow(r.Context(), req.FollowerID, followeeID)
		if err != nil {
			log.Println(err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
