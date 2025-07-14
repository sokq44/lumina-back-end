package handlers

import (
	"backend/config"
	"backend/models"
	"backend/utils/crypt"
	"backend/utils/problems"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func VerifyEmail(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		Token string `json:"token"`
	}

	var body RequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		p := problems.Problem{
			Type:          problems.HandlerProblem,
			ServerMessage: fmt.Sprintf("error while retrieving the access_token cookie: %v", err),
			ClientMessage: "An unexpected error has occurred while processing your request.",
			Status:        http.StatusBadRequest,
		}
		p.Handle(w, r)
		return
	}

	emailValidation, p := db.GetEmailVerificationByToken(body.Token)
	if p.Handle(w, r) {
		return
	}

	if emailValidation.Expires.Before(time.Now()) {
		p := problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: "email validation token has expired",
			ClientMessage: "The verification link is invalid or has expired.",
			Status:        http.StatusGone,
		}
		p.Handle(w, r)
		return
	}

	if db.DeleteEmailVerificationById(emailValidation.Id).Handle(w, r) {
		return
	}

	if db.VerifyUser(emailValidation.UserId).Handle(w, r) {
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	user, p := GetUserFromRequest(r)
	if p.Handle(w, r) {
		return
	}

	userData := map[string]string{
		"id":       user.Id,
		"username": user.Username,
		"email":    user.Email,
		"image":    user.ImageUrl,
	}
	if err := json.NewEncoder(w).Encode(userData); err != nil {
		p := problems.Problem{
			Type:          problems.HandlerProblem,
			ServerMessage: fmt.Sprintf("error while encoding json data to the response: %v", err),
			ClientMessage: "An unexpected error has occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
		p.Handle(w, r)
		return
	}
}

func ModifyUser(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		ImageUrl string `json:"image"`
	}

	var body RequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		p := problems.Problem{
			Type:          problems.HandlerProblem,
			ServerMessage: fmt.Sprintf("error while retrieving the access_token cookie: %v", err),
			ClientMessage: "An unexpected error has occurred while processing your request.",
			Status:        http.StatusBadRequest,
		}
		p.Handle(w, r)
		return
	}

	user, p := db.GetUserByEmail(body.Email)
	if p.Handle(w, r) {
		return
	}

	var newUser = models.User{
		Id:       user.Id,
		Username: body.Username,
		Email:    body.Email,
		ImageUrl: body.ImageUrl,
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
		e := problems.Problem{
			Type:          problems.HandlerProblem,
			ServerMessage: fmt.Sprintf("error while retrieving the access_token cookie: %v", err),
			ClientMessage: "An unexpected error has occurred while processing your request.",
			Status:        http.StatusBadRequest,
		}
		if e.Handle(w, r) {
			return
		}
	}

	user, p := db.GetUserByEmail(body.Email)
	if p.Handle(w, r) {
		return
	}

	token, p := crypt.RandomString(128)
	if p.Handle(w, r) {
		return
	}

	duration := time.Duration(config.PasswdChangeTime)
	passwordChange := models.PasswordChange{
		Token:   token,
		UserId:  user.Id,
		Expires: time.Now().Add(duration),
	}
	if db.CreatePasswordChange(passwordChange).Handle(w, r) {
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
		e := problems.Problem{
			Type:          problems.HandlerProblem,
			ServerMessage: fmt.Sprintf("error while retrieving the access_token cookie: %v", err),
			ClientMessage: "An unexpected error has occurred while processing your request.",
			Status:        http.StatusBadRequest,
		}
		if e.Handle(w, r) {
			return
		}
	}

	passwordChange, p := db.GetPasswordChangeByToken(body.Token)
	if p.Handle(w, r) {
		return
	}

	user, p := db.GetUserById(passwordChange.UserId)
	if p.Handle(w, r) {
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
	db.UpdateUser(*user).Handle(w, r)
}

func PasswordChangeValid(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	token := query.Get("token")

	passwordChange, p := db.GetPasswordChangeByToken(token)
	if p.Handle(w, r) {
		return
	}

	if passwordChange.Expires.Before(time.Now()) {
		p := problems.Problem{
			Type:          problems.HandlerProblem,
			ServerMessage: "password change token has expired.",
			ClientMessage: "The link you used has expired. For security purposes, password reset links are only valid for a limited time. Please request a new link to reset your password.",
			Status:        http.StatusGone,
		}
		p.Handle(w, r)
	}
}

func EmailChangeInit(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		NewEmail string `json:"new_email"`
	}

	var body RequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		p := problems.Problem{
			Type:          problems.HandlerProblem,
			ServerMessage: fmt.Sprintf("error while decoding the request body: %v", err),
			ClientMessage: "An unexpected error has occurred while processing your request.",
			Status:        http.StatusBadRequest,
		}
		p.Handle(w, r)
		return
	}

	u, _ := db.GetUserByEmail(body.NewEmail)
	if u != nil {
		p := &problems.Problem{
			Type:          problems.HandlerProblem,
			ServerMessage: fmt.Sprintf("user with an email address of %s already exists", body.NewEmail),
			ClientMessage: "A user with this email address already exists.",
			Status:        http.StatusConflict,
		}
		p.Handle(w, r)
		return
	}

	user, p := GetUserFromRequest(r)
	if p.Handle(w, r) {
		return
	}

	token, p := crypt.RandomString(128)
	if p.Handle(w, r) {
		return
	}

	duration := time.Duration(config.EmailChangeTime)
	emailChange := models.EmailChange{
		Token:    token,
		NewEmail: body.NewEmail,
		UserId:   user.Id,
		Expires:  time.Now().Add(duration),
	}
	if db.CreateEmailChange(emailChange).Handle(w, r) {
		return
	}

	if em.SendEmailChangeEmail(body.NewEmail, token).Handle(w, r) {
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func ChangeEmail(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		Token string `json:"token"`
	}

	var body RequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		p := problems.Problem{
			Type:          problems.HandlerProblem,
			ServerMessage: fmt.Sprintf("error while decoding the request body: %v", err),
			ClientMessage: "An unexpected error has occurred while processing your request.",
			Status:        http.StatusBadRequest,
		}
		p.Handle(w, r)
		return
	}

	emailChange, p := db.GetEmailChangeByToken(body.Token)
	if p.Handle(w, r) {
		return
	}

	if emailChange.Expires.Before(time.Now()) {
		p := problems.Problem{
			Type:          problems.HandlerProblem,
			ServerMessage: "email change token has expired.",
			ClientMessage: "The link you used has expired. For security purposes, email change links are only valid for a limited time. Please request a new link to change your email.",
			Status:        http.StatusGone,
		}
		p.Handle(w, r)
		return
	}

	user, p := db.GetUserById(emailChange.UserId)
	if p.Handle(w, r) {
		return
	}

	user.Email = emailChange.NewEmail
	if db.UpdateUser(*user).Handle(w, r) {
		return
	}

	db.DeleteEmailChangeById(emailChange.Id).Handle(w, r)
}
