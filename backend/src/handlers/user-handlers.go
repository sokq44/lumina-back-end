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

// not allowed
// conflict
// bad request
// internal server error
// created
func RegisterUserHandler(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		responseWriter.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var u models.User

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

func VerifyEmailHandler(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPatch {
		responseWriter.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var token string
	if err := json.NewDecoder(request.Body).Decode(&token); err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	userId, expires, err := utils.Db.GetEmailValidation(token)
	if err != nil {
		log.Println(err.Error())
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println(userId, expires)
	responseWriter.WriteHeader(http.StatusNoContent)
}
