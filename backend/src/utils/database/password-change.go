package database

import (
	"backend/models"
	"backend/utils/errhandle"
	"database/sql"
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
			Type:          errhandle.DatabaseError,
			ServerMessage: fmt.Sprintf("error while creating a password change token: %v", err),
			ClientMessage: "There was an error while trying to initialize a password change procedure.",
			Status:        http.StatusInternalServerError,
		}
	}

	return nil
}

func (db *Database) GetPasswordChangeByToken(token string) (*models.PasswordChange, *errhandle.Error) {
	var id string
	var tk string
	var userId string
	var expires string

	err := db.Connection.QueryRow(
		"SELECT id, token, expires, user_id FROM password_change WHERE token=?;",
		token,
	).Scan(&id, &tk, &expires, &userId)

	if err == sql.ErrNoRows {
		return nil, &errhandle.Error{
			Type:          errhandle.DatabaseError,
			ServerMessage: fmt.Sprintf("error while getting an email verification by token: %v", err),
			ClientMessage: "This URL is invalid or has expired.",
			Status:        http.StatusNotFound,
		}
	} else if err != nil {
		return nil, &errhandle.Error{
			Type:          errhandle.DatabaseError,
			ServerMessage: fmt.Sprintf("error while getting an email verification by token: %v", err),
			Status:        http.StatusInternalServerError,
		}
	}

	expiresTime, e := parseTime(expires)
	if e != nil {
		return nil, e
	}

	passwordChange := &models.PasswordChange{
		Id:      id,
		Token:   tk,
		UserId:  userId,
		Expires: expiresTime,
	}

	return passwordChange, nil
}

func (db *Database) DeletePasswordChangeById(id string) *errhandle.Error {
	_, err := db.Connection.Exec(
		"DELETE FROM password_change WHERE id=?;",
		id,
	)

	if err != nil {
		return &errhandle.Error{
			Type:          errhandle.DatabaseError,
			ServerMessage: fmt.Sprintf("error while deleting a password change by id: %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
	}

	return nil
}
