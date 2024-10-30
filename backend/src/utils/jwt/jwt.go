package jwt

import (
	"backend/config"
	"backend/models"
	"backend/utils/cryptography"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
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

func GenerateToken(claims Claims) (string, error) {
	header, err := CreateHeader()
	if err != nil {
		return "", err
	}

	payload, err := CreatePayload(claims)
	if err != nil {
		return "", err
	}

	secret := config.JwtSecret

	headerPayload := fmt.Sprintf("%s.%s", header, payload)
	signature := CreateSignature(headerPayload, secret)
	newJWT := fmt.Sprintf("%s.%s.%s", header, payload, signature)

	return newJWT, nil
}

func GenerateAccessToken(user models.User, now time.Time) (string, error) {
	expires := time.Duration(config.JwtAccExpTime)
	claims := Claims{
		"user": user.Id,
		"exp":  now.Add(expires).Unix(),
		"iat":  now.Unix(),
	}

	token, err := GenerateToken(claims)
	if err != nil {
		return "", err
	}

	return token, nil
}

func GenerateRefreshToken(userId string, now time.Time) (models.RefreshToken, error) {
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

func GetRefAccFromRequest(r *http.Request) (string, string, error) {
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
