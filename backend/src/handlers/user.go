package handlers

import (
	"backend/config"
	"backend/models"
	"backend/utils/crypt"
	"backend/utils/database"
	"backend/utils/emails"
	"backend/utils/jwt"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

var db *database.Database = database.GetDb()
var em *emails.SmtpClient = emails.GetEmails()

// TODO: Implement some kind of verification whether the sent data is valid
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
		Password: crypt.Sha256(r.Password),
	}

	exists, e := db.UserExists(u)
	if e.Handle(responseWriter) {
		return
	}
	if exists {
		responseWriter.WriteHeader(http.StatusConflict)
		return
	}

	userId, e := db.CreateUser(u)
	if e.Handle(responseWriter) {
		return
	}

	token, e := crypt.RandomString(128)
	if e.Handle(responseWriter) {
		return
	}

	duration := time.Duration(config.EmailVerTime)
	verification := models.EmailVerification{
		Token:   token,
		UserId:  userId,
		Expires: time.Now().Add(duration),
	}
	e = db.CreateEmailVerification(verification)
	if e.Handle(responseWriter) {
		return
	}

	e = em.SendVerificationEmail(u.Email, token)
	if e.Handle(responseWriter) {
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

	emailValidation, e := db.GetEmailVerificationByToken(body.Token)
	if e.Handle(responseWriter) {
		return
	}

	if emailValidation.Expires.Before(time.Now()) {
		log.Println("token expired")
		responseWriter.WriteHeader(http.StatusGone)
		return
	}

	e = db.DeleteEmailVerificationById(emailValidation.Id)
	if e.Handle(responseWriter) {
		return
	}

	e = db.VerifyUser(emailValidation.UserId)
	if e.Handle(responseWriter) {
		return
	}

	responseWriter.WriteHeader(http.StatusNoContent)
}

var UserLoggedIn http.HandlerFunc = func(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		responseWriter.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

var LoginUser http.HandlerFunc = func(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		responseWriter.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	type RequestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var r RequestBody
	if err := json.NewDecoder(request.Body).Decode(&r); err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, e := db.GetUserByEmail(r.Email)
	if e.Handle(responseWriter) {
		return
	}

	refreshToken, e := db.GetRefreshTokenByUserId(user.Id)
	if refreshToken != nil {
		responseWriter.WriteHeader(http.StatusOK)
		return
	} else if e.Handle(responseWriter) {
		return
	}

	hashedPasswd := crypt.Sha256(r.Password)
	if !user.Verified || hashedPasswd != user.Password {
		responseWriter.WriteHeader(http.StatusForbidden)
		return
	}

	now := time.Now()
	access, e := jwt.GenerateAccessToken(user, now)
	if e.Handle(responseWriter) {
		return
	}

	refresh, e := jwt.GenerateRefreshToken(user.Id, now)
	if e.Handle(responseWriter) {
		return
	}

	e = db.CreateRefreshToken(refresh)
	if e.Handle(responseWriter) {
		return
	}

	http.SetCookie(responseWriter, &http.Cookie{
		Name:     "access_token",
		Value:    access,
		HttpOnly: true,
		Path:     "/",
		Expires:  now.Add(time.Duration(config.JwtAccExpTime)),
	})

	http.SetCookie(responseWriter, &http.Cookie{
		Name:     "refresh_token",
		Value:    refresh.Token,
		HttpOnly: true,
		Path:     "/",
		Expires:  now.Add(time.Duration(config.JwtRefExpTime)),
	})

	responseWriter.WriteHeader(http.StatusOK)
}

var LogoutUser http.HandlerFunc = func(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodDelete {
		responseWriter.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	refreshCookie, err := request.Cookie("refresh_token")
	if err != nil {
		log.Println(err)
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	e := db.DeleteRefreshTokenByToken(refreshCookie.Value)
	if e.Handle(responseWriter) {
		return
	}

	http.SetCookie(responseWriter, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Unix(0, 0),
	})

	http.SetCookie(responseWriter, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Unix(0, 0),
	})

	responseWriter.WriteHeader(http.StatusOK)
}

var GetUser http.HandlerFunc = func(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		responseWriter.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	access, err := request.Cookie("access_token")
	if err != nil {
		log.Println(err)
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	claims, e := jwt.DecodePayload(access.Value)
	if e.Handle(responseWriter) {
		return
	}

	userId := claims["user"].(string)
	user, e := db.GetUserById(userId)
	if e.Handle(responseWriter) {
		return
	}

	userData := map[string]string{
		"username": user.Username,
		"email":    user.Email,
	}
	if err := json.NewEncoder(responseWriter).Encode(userData); err != nil {
		log.Println(err)
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// TODO: Implement some kind of verification whether the sent data is valid
var ModifyUser http.HandlerFunc = func(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPatch {
		responseWriter.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	type RequestBody struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	var body RequestBody
	if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
		log.Println(err)
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, e := db.GetUserByEmail(body.Email)
	if e.Handle(responseWriter) {
		return
	}

	var newUser models.User = models.User{
		Id:       user.Id,
		Username: body.Username,
		Email:    body.Email,
		Password: user.Password,
		Verified: user.Verified,
	}
	e = db.UpdateUser(newUser)
	if e.Handle(responseWriter) {
		return
	}
}
