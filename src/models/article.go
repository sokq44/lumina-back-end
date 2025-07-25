package models

import "time"

type Article struct {
	Id        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserId    string    `json:"user_id"`
	BannerUrl string    `json:"banner"`
	Public    bool      `json:"public"`
	CreatedAt time.Time `json:"created_at"`
}
