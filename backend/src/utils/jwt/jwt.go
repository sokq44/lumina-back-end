package jwt

import (
	"backend/config"
	"backend/models"
	"backend/utils/cryptography"
	"backend/utils/database"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Claims map[string]interface{}

func CreateHeader() (string, error) {
	header := map[string]string{
		"typ": "JWT",
		"alg": "HS256",
	}

	headerJson, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("error while trying to create a header for a JWT: %v", err)
	}

	return cryptography.Base64UrlEncode(headerJson), nil
}

func CreatePayload(claims Claims) (string, error) {
	payloadJson, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("error while trying to create a payload for a JWT: %v", err)
	}

	return cryptography.Base64UrlEncode(payloadJson), nil
}

func CreateSignature(headerPayload, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(headerPayload))

	return cryptography.Base64UrlEncode(h.Sum(nil))
}

func Generate(claims Claims) (string, error) {
	header, err := CreateHeader()
	if err != nil {
		return "", err
	}

	payload, err := CreatePayload(claims)
	if err != nil {
		return "", err
	}

	secret := config.Application.JWT_SECRET

	headerPayload := fmt.Sprintf("%s.%s", header, payload)
	signature := CreateSignature(headerPayload, secret)
	newJWT := fmt.Sprintf("%s.%s.%s", header, payload, signature)

	return newJWT, nil
}

func GenerateAccess(user models.User, now time.Time) (string, error) {
	expires := time.Duration(config.Application.JWT_ACCESS_EXP_TIME)
	claims := Claims{
		"user": user.Id,
		"exp":  now.Add(expires).Unix(),
		"iat":  now.Unix(),
	}

	token, err := Generate(claims)
	if err != nil {
		return "", err
	}

	return token, nil
}

func GenerateRefresh(userId string, now time.Time) (models.RefreshToken, error) {
	expires := time.Duration(config.Application.JWT_REFRESH_EXP_TIME)
	id := uuid.New().String()
	claims := Claims{
		"user": userId,
		"exp":  now.Add(expires).Unix(),
		"jti":  id,
	}

	tk, err := Generate(claims)
	if err != nil {
		return models.RefreshToken{}, err
	}

	return models.RefreshToken{
		Id:      id,
		Token:   tk,
		Expires: now.Add(expires),
		UserId:  userId,
	}, nil
}

func DecodePayload(token string) (Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid jwt token")
	}

	payloadPart := parts[1]
	payloadBytes, err := cryptography.Base64UrlDecode(payloadPart)

	if err != nil {
		return nil, err
	}

	var claims Claims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, fmt.Errorf("error while unmarshaling: %v", err)
	}

	return claims, nil
}

func RefreshAndAccessFromRequest(r *http.Request) (string, string, error) {
	access, err := r.Cookie("access_token")
	if err == http.ErrNoCookie {
		return "", "", err
	} else if err != nil {
		return "", "", fmt.Errorf("error while trying to get the access_token cookie: %v", err)
	}

	refresh, err := r.Cookie("refresh_token")
	if err == http.ErrNoCookie {
		return "", "", err
	} else if err != nil {
		return "", "", fmt.Errorf("error while trying to get the refresh_token cookie: %v", err)
	}

	return access.Value, refresh.Value, nil
}

func Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db := database.GetDb()

		log.Println("Hello from middleware")

		/* Check whether refresh token was passed in the request. */
		refreshToken, err := r.Cookie("refresh_token")
		if err == http.ErrNoCookie {
			log.Println("Doesn't exist...")
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if err != nil {
			log.Printf("error while trying to retrieve the refresh token: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		/* Check whether the refresh token has expired. If it has, delete the cookies and reply with 401.*/
		claims, err := DecodePayload((refreshToken.Value))
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
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if tk.UserId != userId {
			log.Println("Id's don't match...")
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

		/* Check whether the access token has expired, if it did issue another one. */
		if accessToken.Expires.Before(now) {
			claims, err := DecodePayload(accessToken.Value)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			expires := time.Duration(config.Application.JWT_ACCESS_EXP_TIME)
			claims["exp"] = now.Add(expires)
			user, err := database.GetDb().GetUserById(claims["user"].(string))
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			access, err := GenerateAccess(user, now)
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
				Expires:  now.Add(time.Duration(config.Application.JWT_ACCESS_EXP_TIME)),
			})

			next.ServeHTTP(w, r)
			return
		}

		next.ServeHTTP(w, r)
	}
}
