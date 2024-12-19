package models

import "time"

type Article struct {
	Id        int
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserId    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}
