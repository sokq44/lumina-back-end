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
	"fmt"
	"net/http"
	"time"
)

var db = database.GetDb()
var em = emails.GetEmails()

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var body RequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		e := errhandle.Error{
			Type:          errhandle.HandlerError,
			ServerMessage: fmt.Sprintf("error while retrieving the access_token cookie: %v", err),
			ClientMessage: "An error has occurred while processing your request.",
			Status:        http.StatusBadRequest,
		}
		if e.Handle(w, r) {
			return
		}
	}

	u := models.User{
		Username: body.Username,
		Email:    body.Email,
		Password: body.Password,
	}
	if u.Validate(false).Handle(w, r) {
		return
	}

	exists, e := db.UserExists(u)
	if e.Handle(w, r) {
		return
	}
	if exists {
		e := errhandle.Error{
			Type:          errhandle.DatabaseError,
			ServerMessage: "user already exists",
			ClientMessage: "A user with these credentials already exists.",
			Status:        http.StatusConflict,
		}
		if e.Handle(w, r) {
			return
		}
	}

	u.Password = crypt.Sha256(body.Password)
	if db.CreateUser(u).Handle(w, r) {
		return
	}

	token, e := crypt.RandomString(128)
	if e.Handle(w, r) {
		return
	}

	duration := time.Duration(config.EmailVerTime)
	verification := models.EmailVerification{
		Token:   token,
		UserId:  u.Id,
		Expires: time.Now().Add(duration),
	}
	if db.CreateEmailVerification(verification).Handle(w, r) {
		return
	}

	if em.SendVerificationEmail(u.Email, token).Handle(w, r) {
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func VerifyEmail(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		Token string `json:"token"`
	}

	var body RequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		e := errhandle.Error{
			Type:          errhandle.HandlerError,
			ServerMessage: fmt.Sprintf("error while retrieving the access_token cookie: %v", err),
			ClientMessage: "An error has occurred while processing your request.",
			Status:        http.StatusBadRequest,
		}
		if e.Handle(w, r) {
			return
		}
	}

	emailValidation, e := db.GetEmailVerificationByToken(body.Token)
	if e.Handle(w, r) {
		return
	}

	if emailValidation.Expires.Before(time.Now()) {
		e := errhandle.Error{
			Type:          errhandle.DatabaseError,
			ServerMessage: "email validation token has expired",
			ClientMessage: "The verification link is invalid or has expired.",
			Status:        http.StatusGone,
		}
		e.Handle(w, r)
		return
	}

	e = db.DeleteEmailVerificationById(emailValidation.Id)
	if e.Handle(w, r) {
		return
	}

	e = db.VerifyUser(emailValidation.UserId)
	if e.Handle(w, r) {
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var body RequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		e := errhandle.Error{
			Type:          errhandle.HandlerError,
			ServerMessage: fmt.Sprintf("error while retrieving the access_token cookie: %v", err),
			ClientMessage: "An error has occurred while processing your request.",
			Status:        http.StatusBadRequest,
		}
		if e.Handle(w, r) {
			return
		}
	}

	user, e := db.GetUserByEmail(body.Email)
	if e.Handle(w, r) {
		return
	}

	refreshToken, _ := db.GetRefreshTokenByUserId(user.Id)
	if refreshToken != nil && time.Now().After(refreshToken.Expires) {
		db.DeleteRefreshTokenById(refreshToken.Id)
	} else if refreshToken != nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	hashedPasswd := crypt.Sha256(body.Password)
	if !user.Verified || hashedPasswd != user.Password {
		e := errhandle.Error{
			Type:          errhandle.HandlerError,
			ServerMessage: "provided password is incorrect or the user isn't verified:",
			ClientMessage: "Provided password is wrong or the specified user isn't verified.",
			Status:        http.StatusUnauthorized,
		}
		if e.Handle(w, r) {
			return
		}
	}

	now := time.Now()
	access, e := jwt.GenerateAccessToken(user.Id, now)
	if e.Handle(w, r) {
		return
	}

	refresh, e := jwt.GenerateRefreshToken(user.Id, now)
	if e.Handle(w, r) {
		return
	}

	e = db.CreateRefreshToken(refresh)
	if e.Handle(w, r) {
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    access,
		HttpOnly: true,
		Path:     "/",
		Expires:  now.Add(time.Duration(config.JwtAccExpTime)),
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refresh.Token,
		HttpOnly: true,
		Path:     "/",
		Expires:  now.Add(time.Duration(config.JwtRefExpTime)),
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})

	w.WriteHeader(http.StatusOK)
}

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	refreshCookie, err := r.Cookie("refresh_token")
	if err != nil {
		e := errhandle.Error{
			Type:          errhandle.HandlerError,
			ServerMessage: fmt.Sprintf("error while retrieving the refresh_token cookie: %v", err),
			ClientMessage: "An error has occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
		if e.Handle(w, r) {
			return
		}
	}

	e := db.DeleteRefreshTokenByToken(refreshCookie.Value)
	if e.Handle(w, r) {
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

func GetUser(w http.ResponseWriter, r *http.Request) {
	access, err := r.Cookie("access_token")
	if err != nil {
		e := errhandle.Error{
			Type:          errhandle.HandlerError,
			ServerMessage: fmt.Sprintf("error while retrieving the access_token cookie: %v", err),
			ClientMessage: "An error has occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
		if e.Handle(w, r) {
			return
		}
	}

	claims, e := jwt.DecodePayload(access.Value)
	if e.Handle(w, r) {
		return
	}

	userId := claims["user"].(string)
	user, e := db.GetUserById(userId)
	if e.Handle(w, r) {
		return
	}

	userData := map[string]string{
		"username": user.Username,
		"email":    user.Email,
	}
	if err := json.NewEncoder(w).Encode(userData); err != nil {
		e := errhandle.Error{
			Type:          errhandle.HandlerError,
			ServerMessage: fmt.Sprintf("error while encoding json data to the response: %v", err),
			ClientMessage: "An error has occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
		if e.Handle(w, r) {
			return
		}
	}
}

func ModifyUser(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	var body RequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		e := errhandle.Error{
			Type:          errhandle.HandlerError,
			ServerMessage: fmt.Sprintf("error while retrieving the access_token cookie: %v", err),
			ClientMessage: "An error has occurred while processing your request.",
			Status:        http.StatusBadRequest,
		}
		if e.Handle(w, r) {
			return
		}
	}

	user, e := db.GetUserByEmail(body.Email)
	if e.Handle(w, r) {
		return
	}

	var newUser = models.User{
		Id:       user.Id,
		Username: body.Username,
		Email:    body.Email,
		Password: user.Password,
		Verified: user.Verified,
	}
	if newUser.Validate(true).Handle(w, r) {
		return
	}
	if db.UpdateUser(newUser).Handle(w, r) {
		return
	}
}

func PasswordChangeInit(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		Email string `json:"email"`
	}

	var body RequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		e := errhandle.Error{
			Type:          errhandle.HandlerError,
			ServerMessage: fmt.Sprintf("error while retrieving the access_token cookie: %v", err),
			ClientMessage: "An error has occurred while processing your request.",
			Status:        http.StatusBadRequest,
		}
		if e.Handle(w, r) {
			return
		}
	}

	u, e := db.GetUserByEmail(body.Email)
	if e.Handle(w, r) {
		return
	}

	token, e := crypt.RandomString(128)
	if e.Handle(w, r) {
		return
	}

	duration := time.Duration(config.PasswdChangeTime)
	passwdChange := models.PasswordChange{
		Token:   token,
		UserId:  u.Id,
		Expires: time.Now().Add(duration),
	}
	if db.CreatePasswordChange(passwdChange).Handle(w, r) {
		return
	}

	if em.SendPasswordChangeEmail(body.Email, token).Handle(w, r) {
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		Password string `json:"password"`
		Token    string `json:"token"`
	}

	var body RequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		e := errhandle.Error{
			Type:          errhandle.HandlerError,
			ServerMessage: fmt.Sprintf("error while retrieving the access_token cookie: %v", err),
			ClientMessage: "An error has occurred while processing your request.",
			Status:        http.StatusBadRequest,
		}
		if e.Handle(w, r) {
			return
		}
	}

	passwordChange, e := db.GetPasswordChangeByToken(body.Token)
	if e.Handle(w, r) {
		return
	}

	user, e := db.GetUserById(passwordChange.UserId)
	if e.Handle(w, r) {
		return
	}

	user.Password = body.Password
	if user.Validate(false).Handle(w, r) {
		return
	}

	if db.DeletePasswordChangeById(passwordChange.Id).Handle(w, r) {
		return
	}

	user.Password = crypt.Sha256(body.Password)
	if db.UpdateUser(*user).Handle(w, r) {
		return
	}

	w.WriteHeader(http.StatusOK)
}
