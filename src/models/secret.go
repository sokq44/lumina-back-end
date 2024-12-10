package models

import "time"

type Secret struct {
	Id      string
	Secret  string
	Expires time.Time
}
