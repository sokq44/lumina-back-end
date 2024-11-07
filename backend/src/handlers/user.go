package handlers

import (
	"backend/config"
	"backend/models"
	"backend/utils/crypt"
	"backend/utils/database"
	"backend/utils/emails"
	"backend/utils/errhandle"
	"backend/utils/jwt"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

var db *database.Database = database.GetDb()
var em *emails.SmtpClient = emails.GetEmails()

var RegisterUser http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var body RequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u := models.User{
		Id:       uuid.New().String(),
		Username: body.Username,
		Email:    body.Email,
		Password: body.Password,
	}
	if u.Validate(false).Handle(w) {
		return
	}

	exists, e := db.UserExists(u)
	if e.Handle(w) {
		return
	}
	if exists {
		w.WriteHeader(http.StatusConflict)
		return
	}

	u.Password = crypt.Sha256(body.Password)
	if db.CreateUser(u).Handle(w) {
		return
	}

	token, e := crypt.RandomString(128)
	if e.Handle(w) {
		return
	}

	duration := time.Duration(config.EmailVerTime)
	verification := models.EmailVerification{
		Token:   token,
		UserId:  u.Id,
		Expires: time.Now().Add(duration),
	}
	if db.CreateEmailVerification(verification).Handle(w) {
		return
	}

	if em.SendVerificationEmail(u.Email, token).Handle(w) {
		return
	}

	w.WriteHeader(http.StatusCreated)
}

var VerifyEmail http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		Token string `json:"token"`
	}

	var body RequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Println("Problem while decoding!")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	emailValidation, e := db.GetEmailVerificationByToken(body.Token)
	if e.Handle(w) {
		return
	}

	if emailValidation.Expires.Before(time.Now()) {
		e := errhandle.Error{
			Type:    errhandle.DatabaseError,
			Message: "email validation token has expired",
			Status:  http.StatusGone,
		}
		e.Handle(w)
		return
	}

	e = db.DeleteEmailVerificationById(emailValidation.Id)
	if e.Handle(w) {
		return
	}

	e = db.VerifyUser(emailValidation.UserId)
	if e.Handle(w) {
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

var LoginUser http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var body RequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Println("a")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, e := db.GetUserByEmail(body.Email)
	if e.Handle(w) {
		return
	}

	refreshToken, e := db.GetRefreshTokenByUserId(user.Id)
	if refreshToken != nil {
		w.WriteHeader(http.StatusOK)
		return
	} else if e.Handle(w) {
		return
	}

	hashedPasswd := crypt.Sha256(body.Password)
	if !user.Verified || hashedPasswd != user.Password {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	now := time.Now()
	access, e := jwt.GenerateAccessToken(user.Id, now)
	if e.Handle(w) {
		return
	}

	refresh, e := jwt.GenerateRefreshToken(user.Id, now)
	if e.Handle(w) {
		return
	}

	e = db.CreateRefreshToken(refresh)
	if e.Handle(w) {
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    access,
		HttpOnly: true,
		Path:     "/",
		Expires:  now.Add(time.Duration(config.JwtAccExpTime)),
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refresh.Token,
		HttpOnly: true,
		Path:     "/",
		Expires:  now.Add(time.Duration(config.JwtRefExpTime)),
	})

	w.WriteHeader(http.StatusOK)
}

var LogoutUser http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	refreshCookie, err := r.Cookie("refresh_token")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	e := db.DeleteRefreshTokenByToken(refreshCookie.Value)
	if e.Handle(w) {
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Unix(0, 0),
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Unix(0, 0),
	})

	w.WriteHeader(http.StatusOK)
}

var GetUser http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	access, err := r.Cookie("access_token")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	claims, e := jwt.DecodePayload(access.Value)
	if e.Handle(w) {
		return
	}

	userId := claims["user"].(string)
	user, e := db.GetUserById(userId)
	if e.Handle(w) {
		return
	}

	userData := map[string]string{
		"username": user.Username,
		"email":    user.Email,
	}
	if err := json.NewEncoder(w).Encode(userData); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

var ModifyUser http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	var body RequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, e := db.GetUserByEmail(body.Email)
	if e.Handle(w) {
		return
	}

	var newUser models.User = models.User{
		Id:       user.Id,
		Username: body.Username,
		Email:    body.Email,
		Password: user.Password,
		Verified: user.Verified,
	}
	if newUser.Validate(true).Handle(w) {
		return
	}
	if db.UpdateUser(newUser).Handle(w) {
		return
	}
}

var PasswordChangeInit http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		Email string `json:"email"`
	}

	var body RequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u, e := db.GetUserByEmail(body.Email)
	if e.Handle(w) {
		return
	}

	token, e := crypt.RandomString(128)
	if e.Handle(w) {
		return
	}

	duration := time.Duration(config.PasswdChangeTime)
	passwdChange := models.PasswordChange{
		Token:   token,
		UserId:  u.Id,
		Expires: time.Now().Add(duration),
	}
	if db.CreatePasswordChange(passwdChange).Handle(w) {
		log.Println("Error while creating a password change row")
		return
	}

	if em.SendPasswordChangeEmail(body.Email, token).Handle(w) {
		return
	}

	w.WriteHeader(http.StatusCreated)
}

var ChangePassword http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		Password string `json:"password"`
		Token    string `json:"token"`
	}

	var body RequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	passwordChange, e := db.GetPasswordChangeByToken(body.Token)
	if e.Handle(w) {
		return
	}

	user, e := db.GetUserById(passwordChange.UserId)
	if e.Handle(w) {
		return
	}

	user.Password = body.Password
	if user.Validate(false).Handle(w) {
		return
	}

	user.Password = crypt.Sha256(body.Password)
	if db.UpdateUser(user).Handle(w) {
		return
	}

	if db.DeletePasswordChangeById(passwordChange.Id).Handle(w) {
		return
	}

	w.WriteHeader(http.StatusOK)
}
