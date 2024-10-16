package handlers

import (
	"backend/utils"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// TODO:
// Implement some email verification.

func userExists(u User, db *sql.DB) (bool, error) {
	var id string

	err := db.QueryRow("SELECT id FROM users WHERE username=? or email=?;", u.Username, u.Email).Scan(&id)

	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		log.Println(err.Error())
		return false, errors.New("error while trying to execute the query for checking whether a user exists")
	}

	return true, nil
}

func RegisterUserHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var u User
	db := utils.Db.Connection

	if err := json.NewDecoder(request.Body).Decode(&u); err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	exists, err := userExists(u, db)

	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	if exists {
		responseWriter.WriteHeader(http.StatusConflict)
		return
	}

	queryResult, err := db.Exec("INSERT INTO users (username, email, password) values (?, ?, ?)", u.Username, u.Email, utils.SHA256(u.Password))

	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	if affectedRows, err := queryResult.RowsAffected(); err != nil {
		log.Println("error while trying to get affected rows")
	} else {
		log.Println("Register: Rows affected:", affectedRows)
	}

	responseWriter.WriteHeader(http.StatusCreated)
}
