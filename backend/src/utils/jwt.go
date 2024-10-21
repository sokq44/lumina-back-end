package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

type JWT struct {
	Header    string
	Payload   string
	Signature string
}

func (jwt *JWT) CreateHeader() (string, error) {
	header := map[string]string{
		"typ": "JWT",
		"alg": "HS256",
	}

	headerJson, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("error while trying to create a header for a JWT: %v", err.Error())
	}

	return Crypto.Base64UrlEncode(headerJson), nil
}

func (jwt *JWT) CreatePayload(claims map[string]interface{}) (string, error) {
	payloadJson, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("error while trying to create a payload for a JWT: %v", err.Error())
	}

	return Crypto.Base64UrlEncode(payloadJson), nil
}

// Signature generates the JWT signature for the header and payload using HMAC SHA256.
// FIXME: Replace native HMAC and SHA256 with custom implementation when available.
func (jwt *JWT) CreateSignature(headerPayload, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(headerPayload))

	return Crypto.Base64UrlEncode(h.Sum(nil))
}

func (jwt *JWT) New(secret string, claims map[string]interface{}) (JWT, error) {
	header, err := jwt.CreateHeader()
	if err != nil {
		return JWT{}, err
	}

	payload, err := jwt.CreatePayload(claims)
	if err != nil {
		return JWT{}, err
	}

	headerPayload := fmt.Sprintf("%s.%s", header, payload)
	signature := jwt.CreateSignature(headerPayload, secret)

	newJWT := JWT{
		Header:    header,
		Payload:   payload,
		Signature: signature,
	}

	return newJWT, nil
}

func (jwt *JWT) ToString() string {
	return fmt.Sprintf("%s.%s.%s", jwt.Header, jwt.Payload, jwt.Signature)
}
