package models

import "time"

type RefreshToken struct {
	Id      string
	Token   string
	UserId  string
	Expires time.Time
}
