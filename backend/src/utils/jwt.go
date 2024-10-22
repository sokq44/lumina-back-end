package utils

import (
	"backend/config"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

type JsonWebToken struct {
	Secret string
}

type Claims map[string]interface{}

func NewJWT() *JsonWebToken {
	jwt := JsonWebToken{
		Secret: config.AppContext["JWT_SECRET"].(string),
	}

	return &jwt
}

func (jwt *JsonWebToken) CreateHeader() (string, error) {
	header := map[string]string{
		"typ": "JWT",
		"alg": "HS256",
	}

	headerJson, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("error while trying to create a header for a JWT: %v", err)
	}

	return Crypto.Base64UrlEncode(headerJson), nil
}

func (jwt *JsonWebToken) CreatePayload(claims Claims) (string, error) {
	payloadJson, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("error while trying to create a payload for a JWT: %v", err.Error())
	}

	return Crypto.Base64UrlEncode(payloadJson), nil
}

// FIXME: Replace native HMAC-SHA256 with custom implementation when available.
func (jwt *JsonWebToken) CreateSignature(headerPayload, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(headerPayload))

	return Crypto.Base64UrlEncode(h.Sum(nil))
}

func (jwt *JsonWebToken) CreateJWT(claims Claims) (string, error) {
	header, err := jwt.CreateHeader()
	if err != nil {
		return "", err
	}

	payload, err := jwt.CreatePayload(claims)
	if err != nil {
		return "", err
	}

	headerPayload := fmt.Sprintf("%s.%s", header, payload)
	signature := jwt.CreateSignature(headerPayload, jwt.Secret)
	newJWT := fmt.Sprintf("%s.%s.%s", header, payload, signature)

	return newJWT, nil
}
