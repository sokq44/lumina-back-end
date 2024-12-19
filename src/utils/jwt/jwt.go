package jwt

import (
	"backend/config"
	"backend/models"
	"backend/utils/crypt"
	"backend/utils/database"
	"backend/utils/problems"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Claims map[string]interface{}

var db = database.GetDb()

func CreateHeader() (string, *problems.Problem) {
	header := map[string]string{
		"typ": "JWT",
		"alg": "HS256",
	}

	headerJson, err := json.Marshal(header)
	if err != nil {
		return "", &problems.Problem{
			Type:          problems.JwtProblem,
			ServerMessage: fmt.Sprintf("while creating header -> %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
	}

	return crypt.Base64UrlEncode(headerJson), nil
}

func CreatePayload(claims Claims) (string, *problems.Problem) {
	payloadJson, err := json.Marshal(claims)
	if err != nil {
		return "", &problems.Problem{
			Type:          problems.JwtProblem,
			ServerMessage: fmt.Sprintf("while creating payload -> %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
	}

	return crypt.Base64UrlEncode(payloadJson), nil
}

func CreateSignature(headerPayload, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(headerPayload))

	return crypt.Base64UrlEncode(h.Sum(nil))
}

func GenerateToken(claims Claims) (string, *problems.Problem) {
	header, err := CreateHeader()
	if err != nil {
		return "", err
	}

	payload, err := CreatePayload(claims)
	if err != nil {
		return "", err
	}

	latestSecret, err := db.GetLatestSecrets()
	if err != nil {
		return "", err
	}

	headerPayload := fmt.Sprintf("%s.%s", header, payload)
	signature := CreateSignature(headerPayload, latestSecret[0].Secret)
	newJWT := fmt.Sprintf("%s.%s.%s", header, payload, signature)

	return newJWT, nil
}

func GenerateAccessToken(userId string, now time.Time) (string, *problems.Problem) {
	expires := time.Duration(config.JwtAccExpTime)
	claims := Claims{
		"user": userId,
		"exp":  now.Add(expires).Unix(),
		"iat":  now.Unix(),
	}

	token, err := GenerateToken(claims)
	if err != nil {
		return "", err
	}

	return token, nil
}

func GenerateRefreshToken(userId string, now time.Time) (models.RefreshToken, *problems.Problem) {
	expires := time.Duration(config.JwtRefExpTime)
	id := uuid.New().String()
	claims := Claims{
		"user": userId,
		"exp":  now.Add(expires).Unix(),
		"jti":  id,
	}

	tk, err := GenerateToken(claims)
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

func DecodePayload(token string) (Claims, *problems.Problem) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, &problems.Problem{
			Type:          problems.JwtProblem,
			ServerMessage: "token doesn't contain 3 parts",
			ClientMessage: "Error has occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
	}

	payloadPart := parts[1]
	payloadBytes, err := crypt.Base64UrlDecode(payloadPart)

	if err != nil {
		return nil, err
	}

	var claims Claims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, &problems.Problem{
			Type:          problems.JwtProblem,
			ServerMessage: fmt.Sprintf("while decoding payload-> %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
	}

	return claims, nil
}

func WasGeneratedWithSecret(token string, secret string) bool {
	parts := strings.Split(token, ".")
	headerPayload := fmt.Sprintf("%s.%s", parts[0], parts[1])
	signature := CreateSignature(headerPayload, secret)

	return strings.Compare(signature, parts[2]) == 0
}

func GetRefAccFromRequest(r *http.Request) (string, string, *problems.Problem) {
	access, err := r.Cookie("access_token")
	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		return "", "", &problems.Problem{
			Type:          problems.JwtProblem,
			ServerMessage: fmt.Sprintf("while trying to retrieve the access_token cookie -> %v", err),
			ClientMessage: "An error has occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
	}

	refresh, err := r.Cookie("refresh_token")
	if errors.Is(err, http.ErrNoCookie) {
		return "", "", &problems.Problem{
			Type:          problems.JwtProblem,
			ServerMessage: "no refresh_token cookie present",
			ClientMessage: "There was no authentication medium present in the request.",
			Status:        http.StatusUnauthorized,
		}
	} else if err != nil {
		return "", "", &problems.Problem{
			Type:          problems.JwtProblem,
			ServerMessage: fmt.Sprintf("while trying to retrieve the refresh_token cookie -> %v", err),
			ClientMessage: "An error has occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
	}

	if access != nil {
		return access.Value, refresh.Value, nil
	} else {
		return "", refresh.Value, nil
	}
}
