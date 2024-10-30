package middleware

import (
	"backend/config"
	"backend/utils/database"
	"backend/utils/jwt"
	"database/sql"
	"log"
	"net/http"
	"time"
)

func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db := database.GetDb()

		/* Check whether refresh token was passed in the request. */
		refreshToken, err := r.Cookie("refresh_token")
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if err != nil {
			log.Printf("error while trying to retrieve the refresh token: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		/* Check whether the refresh token has expired. If it has, delete the cookies and reply with 401.*/
		claims, err := jwt.DecodePayload((refreshToken.Value))
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		now := time.Now()
		expires := int64(claims["exp"].(float64))
		if expires < now.Unix() {
			w.WriteHeader(http.StatusUnauthorized)

			if err := db.DeleteRefreshTokenByToken(refreshToken.Value); err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
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

			return
		}

		/* Check whether refresh token is assigned to the right person (db). */
		userId := claims["user"].(string)
		tk, err := db.GetRefreshTokenByUserId(userId)
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if tk.UserId != userId {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		/* Check whether access token was passed in the request. */
		accessToken, err := r.Cookie("access_token")
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if err != nil {
			log.Printf("error while trying to retrieve the access token: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		/* Check whether the access token has expired, if it has, issue another one. */
		if accessToken.Expires.Before(now) {
			claims, err := jwt.DecodePayload(accessToken.Value)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			expires := time.Duration(config.JwtAccExpTime)
			claims["exp"] = now.Add(expires)
			user, err := database.GetDb().GetUserById(claims["user"].(string))
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			access, err := jwt.GenerateAccessToken(user, now)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
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

		next.ServeHTTP(w, r)
	}
}
