package handlers

import (
	"backend/config"
	"backend/models"
	"backend/utils/cryptography"
	"backend/utils/database"
	"backend/utils/emails"
	"backend/utils/jwt"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

var db *database.Database = database.GetDb()

func RegisterUser(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		responseWriter.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var u models.User

	if err := json.NewDecoder(request.Body).Decode(&u); err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	exists, err := db.UserExists(u)
	if err != nil {
		log.Println(err)
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	if exists {
		responseWriter.WriteHeader(http.StatusConflict)
		return
	}

	u.Password = cryptography.Sha256(u.Password)
	userId, err := db.CreateUser(u)

	if err != nil {
		log.Println(err)
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	token, err := cryptography.RandomString(128)

	if err != nil {
		log.Println(err)
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	duration := time.Duration(config.Application.EMAIL_VER_TIME)
	verification := models.EmailVerification{
		Token:   token,
		UserId:  userId,
		Expires: time.Now().Add(duration),
	}
	err = db.CreateEmailVerification(verification)
	if err != nil {
		log.Println(err)
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	em := emails.GetEmails()

	err = em.SendVerificationEmail(u.Email, token)
	if err != nil {
		log.Println(err)
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseWriter.WriteHeader(http.StatusCreated)
}

func VerifyEmail(responseWriter http.ResponseWriter, request *http.Request) {
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

	emailValidation, err := db.GetEmailVerificationFromToken(body.Token)
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

	if err = db.DeleteEmailVerification(emailValidation.Id); err != nil {
		log.Println(err.Error())
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = db.VerifyUser(emailValidation.UserId); err != nil {
		log.Println(err.Error())
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseWriter.WriteHeader(http.StatusNoContent)
}

func LoginUser(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		responseWriter.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	type RequestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var body RequestBody
	if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := db.GetUserByEmail(body.Email)
	if err != nil {
		log.Println(err)
		responseWriter.WriteHeader(http.StatusNotFound)
		return
	}

	hashedPasswd := cryptography.Sha256(body.Password)
	if !user.Verified || hashedPasswd != user.Password {
		responseWriter.WriteHeader(http.StatusForbidden)
		return
	}

	expires := time.Duration(config.Application.JWT_EXP_TIME)
	claims := jwt.Claims{
		"user": user.Id,
		"exp":  time.Now().Add(expires).Unix(),
	}

	token, err := jwt.Generate(claims)
	if err != nil {
		log.Println(err)
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseWriter.Header().Set("Authorization", "Bearer "+token)
	responseWriter.WriteHeader(http.StatusOK)
}
