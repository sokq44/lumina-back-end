package models

import (
	"backend/utils/errhandle"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type User struct {
	Id       string
	Username string `json:"username"`
	Email    string `json:"email"`
	ImageUrl string
	Password string `json:"password"`
	Verified bool
}

func (user *User) Validate(passHashed bool) *errhandle.Error {
	if len(user.Username) < 5 || len(user.Username) > 20 {
		return &errhandle.Error{
			Type:          errhandle.ModelError,
			ServerMessage: "username has to be between 5 and 20 characters long",
			ClientMessage: "Username has to be between 5 and 20 characters long.",
			Status:        http.StatusBadRequest,
		}
	}

	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(user.Email) {
		return &errhandle.Error{
			Type:          errhandle.ModelError,
			ServerMessage: "invalid email address",
			ClientMessage: "Invalid email address was provided.",
			Status:        http.StatusBadRequest,
		}
	}

	if !passHashed && (len(user.Password) < 9 || !hasUppercase(user.Password) || !hasDigit(user.Password) || !hasSpecialChar(user.Password)) {
		return &errhandle.Error{
			Type:          errhandle.ModelError,
			ServerMessage: "password must contain a capital letter, a special character, a digit and be at least 9 characters long",
			ClientMessage: "Password must contain a capital letter, a special character, a digit and be at least 9 characters long.",
			Status:        http.StatusBadRequest,
		}
	}

	return nil
}

func hasUppercase(s string) bool {
	for _, c := range s {
		if c >= 'A' && c <= 'Z' {
			return true
		}
	}
	log.Println("no upper")
	return false
}

func hasDigit(s string) bool {
	for _, c := range s {
		if c >= '0' && c <= '9' {
			return true
		}
	}
	log.Println("no digit")
	return false
}

func hasSpecialChar(s string) bool {
	specialChars := "!#$%&'()*+,-./:;<=>?@[]^_{|}~"
	for _, c := range s {
		if strings.ContainsRune(specialChars, c) {
			return true
		}
	}
	log.Println("no special")
	return false
}
