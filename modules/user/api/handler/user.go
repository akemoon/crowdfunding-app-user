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

func GetMe(svc *user.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := httplib.ParseUUIDHeader(r, userIDHeader)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
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

// TODO: implement subscriptions
// GET /me/subscriptions
// GET /users/{id}/subscriptions
