package database

import (
	"backend/models"
	"backend/utils/errhandle"
	"database/sql"
	"fmt"
	"net/http"
)

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
			ServerMessage: fmt.Sprintf("error while getting a password change by token: %v", err),
			ClientMessage: "This URL is invalid or has expired.",
			Status:        http.StatusNotFound,
		}
	} else if err != nil {
		return nil, &errhandle.Error{
			Type:          errhandle.DatabaseError,
			ServerMessage: fmt.Sprintf("error while getting a password change by token: %v", err),
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

func (db *Database) GetExpiredPasswordChanges() ([]models.PasswordChange, *errhandle.Error) {
	rows, err := db.Connection.Query("SELECT * FROM password_change WHERE expires <= NOW();")

	if err != nil {
		return nil, &errhandle.Error{
			Type:          errhandle.DatabaseError,
			ServerMessage: fmt.Sprintf("error while trying to retrieve expired password changes: %v", err),
			Status:        http.StatusInternalServerError,
		}
	}

	var expired []models.PasswordChange
	for rows.Next() {
		var passwordChange models.PasswordChange
		var rawTime string
		if err := rows.Scan(&passwordChange.Id, &passwordChange.Token, &rawTime, &passwordChange.UserId); err != nil {
			return nil, &errhandle.Error{
				Type:          errhandle.DatabaseError,
				ServerMessage: fmt.Sprintf("error while scanning expired password changes: %v", err),
				Status:        http.StatusInternalServerError,
			}
		}

		parsed, err := parseTime(rawTime)
		if err != nil {
			return nil, err
		}

		passwordChange.Expires = parsed
		expired = append(expired, passwordChange)
	}

	return expired, nil
}
