package handlers

import (
	"backend/config"
	"backend/models"
	"backend/utils/crypt"
	"backend/utils/database"
	"backend/utils/emails"
	"backend/utils/jwt"
	"backend/utils/problems"
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
		p := problems.Problem{
			Type:          problems.HandlerProblem,
			ServerMessage: fmt.Sprintf("error while retrieving the access_token cookie: %v", err),
			ClientMessage: "An error has occurred while processing your request.",
			Status:        http.StatusBadRequest,
		}
		p.Handle(w, r)
		return
	}

	u := models.User{
		Username: body.Username,
		Email:    body.Email,
		ImageUrl: config.Host + "/images/default.png",
		Password: body.Password,
	}
	if u.Validate(false).Handle(w, r) {
		return
	}

	exists, p := db.UserExists(u)
	if p.Handle(w, r) {
		return
	}
	if exists {
		p := problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: "user already exists",
			ClientMessage: "A user with these credentials already exists.",
			Status:        http.StatusConflict,
		}
		p.Handle(w, r)
		return
	}

	u.Password = crypt.Sha256(body.Password)
	if db.CreateUser(u).Handle(w, r) {
		return
	}

	token, p := crypt.RandomString(128)
	if p.Handle(w, r) {
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
		p := problems.Problem{
			Type:          problems.HandlerProblem,
			ServerMessage: fmt.Sprintf("error while retrieving the access_token cookie: %v", err),
			ClientMessage: "An error has occurred while processing your request.",
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

func LoginUser(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var body RequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		p := problems.Problem{
			Type:          problems.HandlerProblem,
			ServerMessage: fmt.Sprintf("error while retrieving the access_token cookie: %v", err),
			ClientMessage: "An error has occurred while processing your request.",
			Status:        http.StatusBadRequest,
		}
		p.Handle(w, r)
		return
	}

	user, p := db.GetUserByEmail(body.Email)
	if p.Handle(w, r) {
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
		p := problems.Problem{
			Type:          problems.HandlerProblem,
			ServerMessage: "provided password is incorrect or the user isn't verified:",
			ClientMessage: "Provided password is wrong or the specified user isn't verified.",
			Status:        http.StatusUnauthorized,
		}
		p.Handle(w, r)
		return
	}

	now := time.Now()
	access, p := jwt.GenerateAccessToken(user.Id, now)
	if p.Handle(w, r) {
		return
	}

	refresh, p := jwt.GenerateRefreshToken(user.Id, now)
	if p.Handle(w, r) {
		return
	}

	if db.CreateRefreshToken(refresh).Handle(w, r) {
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
}

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	refreshCookie, err := r.Cookie("refresh_token")
	if err != nil {
		p := problems.Problem{
			Type:          problems.HandlerProblem,
			ServerMessage: fmt.Sprintf("error while retrieving the refresh_token cookie: %v", err),
			ClientMessage: "An error has occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
		p.Handle(w, r)
		return
	}

	if db.DeleteRefreshTokenByToken(refreshCookie.Value).Handle(w, r) {
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Unix(0, 0),
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Unix(0, 0),
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	access, err := r.Cookie("access_token")
	if err != nil {
		p := problems.Problem{
			Type:          problems.HandlerProblem,
			ServerMessage: fmt.Sprintf("error while retrieving the access_token cookie: %v", err),
			ClientMessage: "An error has occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
		p.Handle(w, r)
		return
	}

	claims, p := jwt.DecodePayload(access.Value)
	if p.Handle(w, r) {
		return
	}

	user, p := db.GetUserById(claims["user"].(string))
	if p.Handle(w, r) {
		return
	}

	userData := map[string]string{
		"username": user.Username,
		"email":    user.Email,
		"image":    user.ImageUrl,
	}
	if err := json.NewEncoder(w).Encode(userData); err != nil {
		p := problems.Problem{
			Type:          problems.HandlerProblem,
			ServerMessage: fmt.Sprintf("error while encoding json data to the response: %v", err),
			ClientMessage: "An error has occurred while processing your request.",
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
			ClientMessage: "An error has occurred while processing your request.",
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
			ClientMessage: "An error has occurred while processing your request.",
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
			ClientMessage: "An error has occurred while processing your request.",
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
	if db.UpdateUser(*user).Handle(w, r) {
		return
	}

	access, refresh, p := jwt.GetRefAccFromRequest(r)
	if p.Handle(w, r) {
		return
	}

	if access != "" && refresh != "" {
		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    "",
			HttpOnly: true,
			Path:     "/",
			Expires:  time.Unix(0, 0),
			Secure:   true,
			SameSite: http.SameSiteNoneMode,
		})
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    "",
			HttpOnly: true,
			Path:     "/",
			Expires:  time.Unix(0, 0),
			Secure:   true,
			SameSite: http.SameSiteNoneMode,
		})

		refreshToken, p := db.GetRefreshTokenByUserId(user.Id)
		if p != nil && p.Status == http.StatusNotFound {
			w.WriteHeader(http.StatusOK)
			return
		} else if p != nil && p.Handle(w, r) {
			return
		}

		if db.DeleteRefreshTokenById(refreshToken.Id).Handle(w, r) {
			return
		}
	}
}
