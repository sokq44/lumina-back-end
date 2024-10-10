package handlers

import (
	database "backend/db"
	"encoding/json"
	"log"
	"net/http"
)

type User struct {
	username string `json:username`
	email    string `json:email`
	passowrd string `json:password`
}

func RegisterUser(responseWriter http.ResponseWriter, request *http.Request) {
	var u User

	if err := json.NewDecoder(request.Body).Decode(&u); err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
	} else {
		if err := database.DB.Ping(); err != nil {
			log.Fatal(err)
		}
	}

	responseWriter.WriteHeader(http.StatusCreated)
}
