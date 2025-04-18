package models

import "time"

type Comment struct {
	Id           string
	UserId       string
	Content      string
	CreatedAt    time.Time
	LastModified time.Time
}
