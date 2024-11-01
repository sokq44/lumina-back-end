package middleware

import (
	"backend/config"
	"backend/utils/database"
	"backend/utils/errhandle"
	"backend/utils/jwt"
	"fmt"
	"log"
	"net/http"
	"time"
)

func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db := database.GetDb()
		now := time.Now()

		accessToken, refreshToken, e := GetRefAccFromRequest(r)
		if e.Handle(w) {
			return
		}

		if !jwt.WasGeneratedWithSecret(refreshToken, config.JwtSecret) || !jwt.WasGeneratedWithSecret(accessToken, config.JwtSecret) {
			log.Println("Token wasn't generated here")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claimsRefresh, e := jwt.DecodePayload(refreshToken)
		if e.Handle(w) {
			return
		}

		claimsAccess, e := jwt.DecodePayload(accessToken)
		if e.Handle(w) {
			return
		}

		/* Check whether the refresh token has expired. If it has, delete the cookies and reply with 401.*/
		expiresRefresh := int64(claimsRefresh["exp"].(float64))
		if expiresRefresh < now.Unix() {
			e := db.DeleteRefreshTokenByToken(refreshToken)
			if e.Handle(w) {
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "refresh_token",
				Value:    "",
				HttpOnly: true,
				Path:     "/",
				Expires:  time.Unix(0, 0),
			})

			http.SetCookie(w, &http.Cookie{
				Name:     "access_token",
				Value:    "",
				HttpOnly: true,
				Path:     "/",
				Expires:  time.Unix(0, 0),
			})

			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		/* Check whether refresh token is assigned to the right person (db). */
		userId := claimsRefresh["user"].(string)
		tk, e := db.GetRefreshTokenByUserId(userId)
		if e.Handle(w) {
			return
		}

		if tk.UserId != userId {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		/* Check whether the access token has expired, if it has, issue another one. */
		expiresAccess := int64(claimsAccess["exp"].(float64))
		if expiresAccess < now.Unix() {
			user, e := db.GetUserById(claimsAccess["user"].(string))
			if e.Handle(w) {
				return
			}

			access, e := jwt.GenerateAccessToken(user, now)
			if e.Handle(w) {
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "access_token",
				Value:    access,
				HttpOnly: true,
				Path:     "/",
				Expires:  now.Add(time.Duration(config.JwtAccExpTime)),
			})
		}

		next(w, r)
	}
}

func GetRefAccFromRequest(r *http.Request) (string, string, *errhandle.Error) {
	access, err := r.Cookie("access_token")
	if err == http.ErrNoCookie {
		return "", "", &errhandle.Error{
			Type:    errhandle.JwtError,
			Message: "no access_token cookie present",
			Status:  http.StatusUnauthorized,
		}
	} else if err != nil {
		return "", "", &errhandle.Error{
			Type:    errhandle.JwtError,
			Message: fmt.Sprintf("while trying to retrieve the access_token cookie -> %v", err),
			Status:  http.StatusInternalServerError,
		}
	}

	refresh, err := r.Cookie("refresh_token")
	if err == http.ErrNoCookie {
		return "", "", &errhandle.Error{
			Type:    errhandle.JwtError,
			Message: "no refresh_token cookie present",
			Status:  http.StatusUnauthorized,
		}
	} else if err != nil {
		return "", "", &errhandle.Error{
			Type:    errhandle.JwtError,
			Message: fmt.Sprintf("while trying to retrieve the refresh_token cookie -> %v", err),
			Status:  http.StatusInternalServerError,
		}
	}

	return access.Value, refresh.Value, nil
}
