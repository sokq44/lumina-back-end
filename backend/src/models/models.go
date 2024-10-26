package models

import "time"

type User struct {
	Id       string
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Verified bool
}

type EmailVerification struct {
	Id      string
	Token   string
	UserId  string
	Expires time.Time
}

type RefreshToken struct {
	Id      string
	Token   string
	UserId  string
	Expires time.Time
}
