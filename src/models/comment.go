package models

import "time"

type Comment struct {
	Id           string    `json:"id"`
	UserId       string    `json:"user_id"`
	Content      string    `json:"content"`
	CreatedAt    time.Time `json:"created_at"`
	LastModified time.Time `json:"last_modified"`
}
