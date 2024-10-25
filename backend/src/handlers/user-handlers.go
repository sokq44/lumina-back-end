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

var RegisterUser http.HandlerFunc = func(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		responseWriter.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	type RequestBody struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var r RequestBody
	if err := json.NewDecoder(request.Body).Decode(&r); err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	u := models.User{
		Username: r.Username,
		Email:    r.Email,
		Password: cryptography.Sha256(r.Password),
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

var VerifyEmail http.HandlerFunc = func(responseWriter http.ResponseWriter, request *http.Request) {
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

var LoginUser http.HandlerFunc = func(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		responseWriter.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	log.Println("hello from login handler")

	type RequestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var r RequestBody
	if err := json.NewDecoder(request.Body).Decode(&r); err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := db.GetUserByEmail(r.Email)
	if err != nil {
		log.Println(err)
		responseWriter.WriteHeader(http.StatusNotFound)
		return
	}

	hashedPasswd := cryptography.Sha256(r.Password)
	if !user.Verified || hashedPasswd != user.Password {
		responseWriter.WriteHeader(http.StatusForbidden)
		return
	}

	now := time.Now()

	access, err := jwt.GenerateAccess(user, now)
	if err != nil {
		log.Println(err.Error())
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	refresh, err := jwt.GenerateRefresh(user.Id, now)
	if err != nil {
		log.Println(err.Error())
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := db.CreateRefreshToken(refresh); err != nil {
		log.Println(err.Error())
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(responseWriter, &http.Cookie{
		Name:     "access_token",
		Value:    access,
		HttpOnly: true,
		Path:     "/",
		Expires:  now.Add(time.Duration(config.Application.JWT_ACCESS_EXP_TIME)),
	})

	http.SetCookie(responseWriter, &http.Cookie{
		Name:     "refresh_token",
		Value:    refresh.Token,
		HttpOnly: true,
		Path:     "/",
		Expires:  now.Add(time.Duration(config.Application.JWT_REFRESH_EXP_TIME)),
	})

	responseWriter.WriteHeader(http.StatusOK)
}
