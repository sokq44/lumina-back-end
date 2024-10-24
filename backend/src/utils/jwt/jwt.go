package jwt

import (
	"backend/config"
	"backend/models"
	"backend/utils/cryptography"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
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
		return "", fmt.Errorf("error while trying to create a payload for a JWT: %v", err.Error())
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
