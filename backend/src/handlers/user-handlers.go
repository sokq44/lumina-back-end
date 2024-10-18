package handlers

import (
	"backend/models"
	"backend/utils"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

const EMAIL_VER_TIME = time.Duration(time.Hour * 3)

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

	err = utils.Db.CreateEmailVerification(models.EmailVerification{
		Token:   token,
		UserId:  userId,
		Expires: time.Now().Add(EMAIL_VER_TIME),
	})

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

	type RequestBody struct {
		Token string `json:"token"`
	}

	var body RequestBody
	if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
		log.Println("Problem while decoding!")
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	emailValidation, err := utils.Db.GetEmailVerificationFromToken(body.Token)
	if err != nil {
		log.Println(err.Error())
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	if emailValidation.Expires.Before(time.Now()) {
		log.Println("token expired")
		responseWriter.WriteHeader(http.StatusGone)
		return
	}

	if err = utils.Db.DeleteEmailVerification(emailValidation.Id); err != nil {
		log.Println(err.Error())
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = utils.Db.VerifyUser(emailValidation.UserId); err != nil {
		log.Println(err.Error())
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseWriter.WriteHeader(http.StatusNoContent)
}
