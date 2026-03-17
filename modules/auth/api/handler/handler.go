package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/akemoon/golib/httplib"
	"github.com/akemoon/crowdfunding-app-user/modules/auth/domain"
	"github.com/akemoon/crowdfunding-app-user/modules/auth/metrics"
	"github.com/akemoon/crowdfunding-app-user/modules/auth/service/auth"
)

const (
	userIDHeader   = "X-User-ID"
	userRoleHeader = "X-User-Role"
)

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

func SignIn(svc *auth.Service, m *metrics.AuthMetrics) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req domain.SignInReq

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			m.AuthSignInTotal.WithLabelValues("failure").Inc()
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		resp, err := svc.SignIn(r.Context(), req)
		if err != nil {
			m.AuthSignInTotal.WithLabelValues("failure").Inc()
			log.Printf("service: %s", err)

			status, errResp := httplib.MapErrToHTTP(err, SignInMapRules)
			httplib.WriteJSON(w, status, errResp)
			return
		}

		m.AuthSignInTotal.WithLabelValues("success").Inc()
		httplib.WriteJSON(w, http.StatusOK, resp)
	}
}

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

func CheckAccessToken(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		optional := r.URL.Query().Get("optional") == "true"

		authHeader := r.Header.Get("Authorization")
		if strings.TrimSpace(authHeader) == "" {
			if optional {
				w.WriteHeader(http.StatusOK)
				return
			}
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		claims, err := svc.ValidateAccessToken(authHeader)
		if err != nil {
			log.Printf("token service: %s", err)
			status, resp := httplib.MapErrToHTTP(err, CheckMapRules)
			httplib.WriteJSON(w, status, resp)
			return
		}

		w.Header().Set(userIDHeader, claims.UserID.String())
		w.Header().Set(userRoleHeader, claims.Role)
		w.WriteHeader(http.StatusOK)
	}
}

func Refresh(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req domain.RefreshReq

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		resp, err := svc.Refresh(r.Context(), req)
		if err != nil {
			log.Printf("service: %s", err)

			status, errResp := httplib.MapErrToHTTP(err, RefreshMapRules)
			httplib.WriteJSON(w, status, errResp)
			return
		}

		httplib.WriteJSON(w, http.StatusOK, resp)
	}
}
