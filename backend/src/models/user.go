package models

type User struct {
	Id       string
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Verified bool
}
