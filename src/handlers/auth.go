package handlers

import (
	"backend/config"
	"backend/models"
	"backend/utils/crypt"
	"backend/utils/jwt"
	"backend/utils/problems"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func Register(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var body RequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		p := problems.Problem{
			Type:          problems.HandlerProblem,
			ServerMessage: fmt.Sprintf("error while decoding request body: %v", err),
			ClientMessage: "An unexpected error has occurred while processing your request.",
			Status:        http.StatusBadRequest,
		}
		p.Handle(w, r)
		return
	}

	u := models.User{
		Id:       uuid.New().String(),
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

func Login(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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

func Logout(w http.ResponseWriter, r *http.Request) {
	refreshCookie, err := r.Cookie("refresh_token")
	if err != nil {
		p := problems.Problem{
			Type:          problems.HandlerProblem,
			ServerMessage: fmt.Sprintf("error while retrieving the refresh_token cookie: %v", err),
			ClientMessage: "An unexpected error has occurred while processing your request.",
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
