package handlers

import (
	"backend/models"
	"backend/utils"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// TODO:
// Implement some email verification.

func RegisterUserHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var u models.User
	// smtp := utils.Smtp

	if err := json.NewDecoder(request.Body).Decode(&u); err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	exists, err := utils.Db.UserExists(u)

	if err != nil {
		log.Println(err)
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	if exists {
		responseWriter.WriteHeader(http.StatusConflict)
		return
	}

	userId, err := utils.Db.CreateUser(u)

	if err != nil {
		log.Println(err)
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	token, err := utils.Encryptor.RandomString(128)

	if err != nil {
		log.Println(err)
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	err = utils.Db.CreateEmailValidation(userId, token, time.Now().Add(time.Hour*3))

	if err != nil {
		log.Println(err)
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	err = utils.Smtp.SendVerificationEmail(u.Email, token)
	if err != nil {
		log.Println(err)
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseWriter.WriteHeader(http.StatusCreated)
}
