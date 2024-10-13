package userHandlers

import (
	database "backend/db"
	hashing "backend/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// TODO:
// Hash the user's password.
// Implement some email verification.
// Verify request for any sql injection.
// check whether user with a certain username or email already exists.

func RegisterUserHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var u User
	db, err := database.GetDbConnection()

	if err != nil {
		fmt.Fprintln(responseWriter, "There was a problem with getting database connection.")
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewDecoder(request.Body).Decode(&u); err != nil {
		fmt.Fprintln(responseWriter, "There was a problem with decoding the request body.")
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	query := fmt.Sprintf("INSERT INTO users (username, email, password) values ('%s', '%s', '%s')", u.Username, u.Email, hashing.SHA256(u.Password))
	queryResult, err := db.Exec(query)

	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(responseWriter, "There was a problem with executing a query for regitering the user.")
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
