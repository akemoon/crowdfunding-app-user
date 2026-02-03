package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/akemoon/crowdfunding-app-user/domain"
	"github.com/akemoon/crowdfunding-app-user/service/user"
	"github.com/google/uuid"
)

// @Summary Create a user
// @Description Create a user with the given payload
// @Accept json
// @Produce json
// @Param payload body domain.CreateUserReq true "Create user request payload"
// @Success 201 {object} domain.CreateUserResp "User created"
// @Failure 400 {object} ErrResp "Invalid request"
// @Failure 409 {object} ErrResp "Username conflict"
// @Failure 405 {string} string "Method not allowed"
// @Failure 500 {object} ErrResp "Internal server error"
// @Router /user [post]
func CreateUser(svc *user.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req domain.CreateUserReq

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		resp, err := svc.CreateUser(r.Context(), req)
		if err != nil {
			log.Printf("service: %s", err.Error())
			status, errResp := MapErrToHTTP(err)
			writeWithJson(w, status, errResp)
			return
		}

		writeWithJson(w, http.StatusCreated, resp)
	}
}

// @Summary Get user by ID
// @Description Get user info by user ID
// @Accept json
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Success 200 {object} domain.User "User information"
// @Failure 400 {string} string "Invalid user ID format"
// @Failure 404 {object} ErrResp "User not found"
// @Failure 405 {string} string "Method not allowed"
// @Failure 500 {object} ErrResp "Internal server error"
// @Router /user/id/{id} [get]
func GetUserByID(svc *user.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		idStr := r.PathValue("id")
		id, err := parseUserID(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		user, err := svc.GetUserByID(r.Context(), id)
		if err != nil {
			status, errResp := MapErrToHTTP(err)
			writeWithJson(w, status, errResp)
			return
		}

		writeWithJson(w, http.StatusOK, user)
	}
}

// @Summary Get user by username
// @Description Get user info by username
// @Accept json
// @Produce json
// @Param username path string true "Username"
// @Success 200 {object} domain.User "User information"
// @Failure 404 {object} ErrResp "User not found"
// @Failure 405 {string} string "Method not allowed"
// @Failure 500 {object} ErrResp "Internal server error"
// @Router /user/{username} [get]
func GetUserByUsername(svc *user.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		username := r.PathValue("username")

		user, err := svc.GetUserByUsername(r.Context(), username)
		if err != nil {
			status, errResp := MapErrToHTTP(err)
			writeWithJson(w, status, errResp)
			return
		}

		writeWithJson(w, http.StatusOK, user)
	}
}

// @Summary Get current user
// @Description Get user info by X-User-ID header
// @Accept json
// @Produce json
// @Param X-User-ID header string true "User ID (UUID)"
// @Success 200 {object} domain.User "User information"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {object} ErrResp "User not found"
// @Failure 405 {string} string "Method not allowed"
// @Failure 500 {object} ErrResp "Internal server error"
// @Router /user/me [get]
func GetMe(svc *user.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		idStr := r.Header.Get("X-User-ID")
		if idStr == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		id, err := parseUserID(idStr)

		user, err := svc.GetUserByID(r.Context(), id)
		if err != nil {
			status, errResp := MapErrToHTTP(err)
			writeWithJson(w, status, errResp)
			return
		}

		writeWithJson(w, http.StatusOK, user)
	}
}

func parseUserID(idStr string) (uuid.UUID, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user id %s", idStr)
	}

	return id, nil
}

func writeWithJson(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
