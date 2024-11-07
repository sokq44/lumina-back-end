package database

import (
	"backend/models"
	"backend/utils/errhandle"
	"fmt"
	"net/http"
)

// TODO: clean out the expired tokens

func (db *Database) CreatePasswordChange(p models.PasswordChange) *errhandle.Error {
	_, err := db.Connection.Exec(
		"INSERT INTO password_change (token, expires, user_id) values (?, ?, ?);",
		p.Token, p.Expires, p.UserId,
	)

	if err != nil {
		return &errhandle.Error{
			Type:    errhandle.DatabaseError,
			Message: fmt.Sprintf("error while creating a password change token: %v", err),
			Status:  http.StatusInternalServerError,
		}
	}

	return nil
}
