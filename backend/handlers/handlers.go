package handlers

import (
	database "backend/db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func RegisterUser(responseWriter http.ResponseWriter, request *http.Request) {
	var u User
	db, err := database.GetDbConnection()

	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
	}

	if err := json.NewDecoder(request.Body).Decode(&u); err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	} else {
		query := fmt.Sprintf("INSERT INTO users (username, email, password) values ('%s', '%s', '%s')", u.Username, u.Email, u.Password)
		_, err := db.Query(query)

		if err != nil {
			log.Println(err.Error())
		} else {
			log.Println("Added user", u.Username)
		}
	}

	responseWriter.WriteHeader(http.StatusCreated)
}
