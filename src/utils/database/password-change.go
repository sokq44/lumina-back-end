package database

import (
	"backend/models"
	"backend/utils/problems"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-sql-driver/mysql"
)

func (db *Database) CreatePasswordChange(p models.PasswordChange) *problems.Problem {
	row := db.Connection.QueryRow("SELECT id, expires FROM password_change WHERE user_id=?", p.UserId)
	if row != nil {
		var id string
		var raw string

		err := row.Scan(&id, &raw)
		if !errors.Is(err, sql.ErrNoRows) {
			if err != nil {
				return &problems.Problem{
					Type:          problems.DatabaseProblem,
					ServerMessage: fmt.Sprintf("while scanning a row for the CreatePasswordChange function -> %v", err),
					ClientMessage: "An unexpected error has occurred while processing your request.",
					Status:        http.StatusInternalServerError,
				}
			}
			expires, p := parseTime(raw)
			if p != nil {
				return p
			}

			if expires.Before(time.Now()) {
				p = db.DeletePasswordChangeById(id)
				if p != nil {
					return p
				}
			}
		}

	}

	_, err := db.Connection.Exec(
		"INSERT INTO password_change (token, expires, user_id) values (?, ?, ?);",
		p.Token, p.Expires, p.UserId,
	)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return &problems.Problem{
				Type:          problems.DatabaseProblem,
				ServerMessage: fmt.Sprintf("duplicate entry error while creating a password change token: %v", err),
				ClientMessage: "A password change request already exists for this user.",
				Status:        http.StatusConflict,
			}
		}
		return &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("error while creating a password change token: %v", err),
			ClientMessage: "There was an error while trying to initialize a password change procedure.",
			Status:        http.StatusInternalServerError,
		}
	}

	return nil
}

func (db *Database) GetPasswordChangeByToken(token string) (*models.PasswordChange, *problems.Problem) {
	var id string
	var tk string
	var userId string
	var expires string

	err := db.Connection.QueryRow(
		"SELECT id, token, expires, user_id FROM password_change WHERE token=?;",
		token,
	).Scan(&id, &tk, &expires, &userId)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("error while getting a password change by token: %v", err),
			ClientMessage: "We couldn't find a valid password reset request. Please check your link or request a new password reset.",
			Status:        http.StatusNotFound,
		}
	} else if err != nil {
		return nil, &problems.Problem{
			Type:          problems.DatabaseProblem,
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

func (db *Database) DeletePasswordChangeById(id string) *problems.Problem {
	_, err := db.Connection.Exec(
		"DELETE FROM password_change WHERE id=?;",
		id,
	)

	if err != nil {
		return &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("error while deleting a password change by id: %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
	}

	return nil
}

func (db *Database) GetExpiredPasswordChanges() ([]models.PasswordChange, *problems.Problem) {
	rows, err := db.Connection.Query("SELECT * FROM password_change WHERE expires <= NOW();")

	if err != nil {
		return nil, &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("error while trying to retrieve expired password changes: %v", err),
			Status:        http.StatusInternalServerError,
		}
	}

	var expired []models.PasswordChange
	for rows.Next() {
		var passwordChange models.PasswordChange
		var rawTime string
		if err := rows.Scan(&passwordChange.Id, &passwordChange.Token, &rawTime, &passwordChange.UserId); err != nil {
			return nil, &problems.Problem{
				Type:          problems.DatabaseProblem,
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
