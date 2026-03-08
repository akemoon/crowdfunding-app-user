package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/akemoon/crowdfunding-app-user/lib/httplib"
	"github.com/akemoon/crowdfunding-app-user/modules/auth/domain"
	"github.com/akemoon/crowdfunding-app-user/modules/auth/service/auth"
)

// @Summary Sign up
// @Description Sign up by given payload
// @Accept json
// @Produce json
// @Param payload body domain.SignUpReq true "Sign up payload"
// @Success 201 "User created"
// @Failure 400 {object} httplib.ErrResp "Validation error"
// @Failure 405 "Method not allowed"
// @Failure 409 {object} httplib.ErrResp "Username or email exists"
// @Failure 500 {object} httplib.ErrResp "Internal server error"
// @Router /signup [post]
func SignUp(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req domain.SignUpReq

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		err = svc.SignUp(r.Context(), req)
		if err != nil {
			log.Printf("service: %s", err.Error())

			status, resp := httplib.MapErrToHTTP(err, SignUpMapRules)
			httplib.WriteJSON(w, status, resp)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

// @Summary Sign in
// @Description Sign in by given paload
// @Accept json
// @Produce json
// @Param payload body domain.SignInReq true "Sign in payload"
// @Success 200 {object} domain.SignInResp "Tokens issued"
// @Failure 401 {object} httplib.ErrResp "Invalid credentials"
// @Failure 405 "Method not allowed"
// @Failure 500 {object} httplib.ErrResp "Internal server error"
// @Router /signin [post]
func SignIn(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req domain.SignInReq

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		resp, err := svc.SignIn(r.Context(), req)
		if err != nil {
			log.Printf("service: %s", err)

			status, errResp := httplib.MapErrToHTTP(err, SignInMapRules)
			httplib.WriteJSON(w, status, errResp)
			return
		}

		httplib.WriteJSON(w, http.StatusOK, resp)
	}
}

// @Summary Sign out
// @Description Revoke refresh token
// @Accept json
// @Produce json
// @Param payload body domain.SignOutReq true "Sign out payload"
// @Success 200 "Signed out"
// @Failure 400 {object} httplib.ErrResp "Invalid refresh token"
// @Failure 405 "Method not allowed"
// @Failure 500 {object} httplib.ErrResp "Internal server error"
// @Router /signout [post]
func SignOut(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req domain.SignOutReq

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		err = svc.SignOut(r.Context(), req)
		if err != nil {
			log.Printf("service: %s", err)

			status, errResp := httplib.MapErrToHTTP(err, SignOutMapRules)
			httplib.WriteJSON(w, status, errResp)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// @#Summary Check access token
// @#Description Validate access token from Authorization header
// @#Accept json
// @#Produce json
// @#Param Authorization header string true "Authorization header with access token"
// @#Success 200 "Access token is valid"
// @#Header  200 {string} X-User-Id "Authenticated user UUID"
// @#Failure 401 "Unauthorized"
// @#Failure 405 "Method not allowed"
// @#Failure 500 "Internal server error"
// @#Router /check [get]
// func CheckAccessToken(svc *token.Service) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method != http.MethodGet {
// 			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
// 			return
// 		}
//
// 		authHeader := r.Header.Get("Authorization")
// 		if strings.TrimSpace(authHeader) == "" {
// 			http.Error(w, "missing authorization header", http.StatusUnauthorized)
// 			return
// 		}
//
// 		userID, err := svc.ValidateAccessToken(authHeader)
// 		if err != nil {
// 			log.Printf("token service: %s", err)
//
// 			status, resp := mapErrToHTTP(err)
// 			writeJSON(w, status, resp)
// 			return
// 		}
//
// 		w.Header().Set("X-User-Id", userID.String())
// 		w.WriteHeader(http.StatusOK)
// 	}
// }