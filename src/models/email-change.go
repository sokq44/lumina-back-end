package models

import "time"

type EmailChange struct {
	Id       string
	Token    string
	NewEmail string
	UserId   string
	Expires  time.Time
}
