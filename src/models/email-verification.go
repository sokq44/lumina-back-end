package models

import "time"

type EmailVerification struct {
	Id      string
	Token   string
	UserId  string
	Expires time.Time
}
