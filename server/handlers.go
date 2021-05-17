package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	COOKIE_CSRF           = "_iotdash_csrf_token"
	COOKIE_SESSION        = "_iotdash_session_id"
	COOKIE_SESSION_STATUS = "_iotdash_session_status"
	HEADER_CSRF           = "iotdash-csrf-token"
	HEADER_BEARER_TOKEN   = "authorization"

	BearerPrefix = "Bearer "
)

func HandleCsrf(verify bool, next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctxAuth := req.Context().Value(AuthContext{})
		authContext, ok := ctxAuth.(AuthContext)
		if !ok || authContext.Method != AuthMethodBasic {
			next.ServeHTTP(w, req)
			return
		}

		csrfTokenCookie := &http.Cookie{
			Name:    COOKIE_CSRF,
			Value:   GetSecureToken(),
			Expires: time.Now().Add(1000000 * time.Minute),
		}
		http.SetCookie(w, csrfTokenCookie)

		if verify {
			csrfHeader := req.Header.Get(HEADER_CSRF)
			csrfCookie, err := req.Cookie(COOKIE_CSRF)

			authorized := false

			if err != nil || csrfHeader == "" {
				authorized = false
			} else if csrfCookie.Value == csrfHeader {
				authorized = true
			}

			if !authorized {
				http.Error(w, "CSRF Token is invalid!", http.StatusUnauthorized)
				return
			}
		}

		next.ServeHTTP(w, req)
	})
}

func HandleSessionAuth(app App, next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		sessionCookie, err := req.Cookie(COOKIE_SESSION)
		if err != nil {
			next.ServeHTTP(w, req)
			return
		}

		sesh := app.SessionGetById(sessionCookie.Value)
		if sesh == nil {
			next.ServeHTTP(w, req)
			return
		}

		acct, _ := app.AccountGetByUsername(sesh.Username)

		ctx := context.WithValue(req.Context(), AuthContext{}, AuthContext{Method: AuthMethodBasic, Account: acct, Session: sesh})
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}

func HandleBearerAuth(app App, next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		bearerHeader := req.Header.Get(HEADER_BEARER_TOKEN)

		token := strings.TrimPrefix(bearerHeader, BearerPrefix)

		acct, err := app.AccountGetByBearerToken(token)

		if acct != nil {
			ctx := context.WithValue(req.Context(), AuthContext{}, AuthContext{Method: AuthMethodBearer, Account: acct})
			req = req.WithContext(ctx)
		} else {
			fmt.Println(err)
		}

		next.ServeHTTP(w, req)
	})
}

func HandleAuth(app App, next http.Handler) http.HandlerFunc {
	return HandleSessionAuth(
		app,
		HandleBearerAuth(
			app,
			HandleCsrf(true, next),
		),
	)
}
