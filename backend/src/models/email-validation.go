package models

import "time"

type EmailValidation struct {
	Id      string
	Token   string
	UserId  string
	Expires time.Time
}
