package jwt

import (
	"backend/config"
	"backend/utils/cryptography"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
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
