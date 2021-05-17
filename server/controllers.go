package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func HandleSecurityLogin(app App) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				http.Error(w, "Unauthorized!", http.StatusUnauthorized)
				fmt.Println(r)
			}
		}()

		dto := LoginPostDTO{}
		err := json.NewDecoder(req.Body).Decode(&dto)
		if err != nil {
			panic("Error deserializing DTO")
		}

		acct, err := app.AccountGetByUsername(dto.Username)
		if err != nil {
			panic("Error looking up user [" + dto.Username + "]: " + err.Error())
		}

		err = bcrypt.CompareHashAndPassword([]byte(acct.AdminPassword), []byte(dto.Password))
		if err != nil {
			panic("Password does not match hash!")
		} else {
			sesh := app.SessionCreate(acct.AdminUsername)

			setSessionCookies(w, sesh)
		}
	})
}

// Deletes user's session if it exists and clears session cookies on the client.
func HandleSecurityLogout(app App) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authCtx, ok := req.Context().Value(AuthContext{}).(AuthContext)
		if !ok || authCtx.Method == AuthMethodBasic {
			app.SessionDelete(authCtx.Session.Id)
		}

		clearSessionCookies(w)
	})
}

func HandleMetrics(app App) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// TODO This should be handled consistently across all API endpoints
		w.Header().Add("Content-Type", "application/json")

		authCtx, ok := req.Context().Value(AuthContext{}).(AuthContext)
		if !ok || authCtx.Method == AuthMethodNone {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		switch req.Method {
		// GET
		case http.MethodGet:
			count, err := app.AccountGetIotUserCount(authCtx.Id)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			dto := MetricsGetResDTO{IotUserCount: count}
			resBody, _ := json.Marshal(dto)
			w.Write(resBody)
			return
		// POST
		case http.MethodPost:
			dto := MetricsPostDTO{}
			err := json.NewDecoder(req.Body).Decode(&dto)
			if err != nil || dto.AccountId == "" || dto.UserId == "" {
				// Or unprocessible entity
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if dto.AccountId != authCtx.Id {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			if err := app.AccountRegisterIotUser(dto.AccountId, dto.UserId); err != nil {
				if err.Error() == "limit reached" {
					fmt.Println("Limit reached")
					// TODO Return an appropriate status code
				} else {
					panic(err)
				}
			}
		}
	})
}

func HandleAccountUpgrade(app App) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// TODO move to a middleware or at least util method
		authCtx, ok := req.Context().Value(AuthContext{}).(AuthContext)
		if !ok || authCtx.Method == AuthMethodNone {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		if req.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

		if err := app.AccountUpgrade(authCtx.Id); err != nil {
			panic(err)
		}
	})
}

// Sets session cookies expiry in the past to clear them from the client.
func clearSessionCookies(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:    COOKIE_SESSION,
		Expires: time.Now().UTC().Add(-1 * time.Hour),
	})

	http.SetCookie(w, &http.Cookie{
		Name:    COOKIE_SESSION_STATUS,
		Expires: time.Now().UTC().Add(-1 * time.Hour),
	})
}

// Sets the session cookies based on the provided session pointer.
func setSessionCookies(w http.ResponseWriter, sesh *Session) {
	// httponly cookie for passing session id to server
	http.SetCookie(w, &http.Cookie{
		Name:     COOKIE_SESSION,
		Value:    sesh.Id,
		Expires:  sesh.Expiry,
		HttpOnly: true,
	})

	// JS accessible cookie to notify UI of session status
	http.SetCookie(w, &http.Cookie{
		Name:    COOKIE_SESSION_STATUS,
		Value:   "true",
		Expires: sesh.Expiry,
	})
}
