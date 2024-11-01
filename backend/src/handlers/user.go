package handlers

import (
	"backend/config"
	"backend/models"
	"backend/utils/crypt"
	"backend/utils/database"
	"backend/utils/emails"
	"backend/utils/jwt"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

var db *database.Database = database.GetDb()

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

	token, err := crypt.RandomString(128)

	if err != nil {
		log.Println(err)
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	duration := time.Duration(config.EmailVerTime)
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

	emailValidation, err := db.GetEmailVerificationByToken(body.Token)
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

	if err = db.DeleteEmailVerificationById(emailValidation.Id); err != nil {
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

	user, err := db.GetUserByEmail(r.Email)
	if err != nil {
		log.Println(err)
		responseWriter.WriteHeader(http.StatusNotFound)
		return
	}

	refreshToken, err := db.GetRefreshTokenByUserId(user.Id)
	if refreshToken != nil {
		responseWriter.WriteHeader(http.StatusOK)
		return
	} else if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	hashedPasswd := crypt.Sha256(r.Password)
	if !user.Verified || hashedPasswd != user.Password {
		responseWriter.WriteHeader(http.StatusForbidden)
		return
	}

	now := time.Now()
	access, err := jwt.GenerateAccessToken(user, now)
	if err != nil {
		log.Println(err.Error())
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	refresh, err := jwt.GenerateRefreshToken(user.Id, now)
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

	if err := db.DeleteRefreshTokenByToken(refreshCookie.Value); err != nil {
		log.Println(err)
		responseWriter.WriteHeader(http.StatusInternalServerError)
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

	claims, err := jwt.DecodePayload(access.Value)
	if err != nil {
		log.Println(err)
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	userId := claims["user"].(string)
	user, err := db.GetUserById(userId)
	if err != nil {
		log.Println(err)
		responseWriter.WriteHeader(http.StatusInternalServerError)
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

	user, err := db.GetUserByEmail(body.Email)
	if err != nil {
		log.Println(err)
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	var newUser models.User = models.User{
		Id:       user.Id,
		Username: body.Username,
		Email:    body.Email,
		Password: user.Password,
		Verified: user.Verified,
	}
	if err := db.UpdateUser(newUser); err != nil {
		log.Println(err)
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
}
