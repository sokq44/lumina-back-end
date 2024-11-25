package models

import "time"

type PasswordChange struct {
	Id      string
	Token   string
	UserId  string
	Expires time.Time
}
