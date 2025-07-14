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

func (db *Database) CreateEmailChange(e models.EmailChange) *problems.Problem {
	row := db.Connection.QueryRow("SELECT id, expires FROM email_change WHERE user_id=?", e.UserId)
	if row != nil {
		var id string
		var raw string

		err := row.Scan(&id, &raw)
		if !errors.Is(err, sql.ErrNoRows) {
			if err != nil {
				return &problems.Problem{
					Type:          problems.DatabaseProblem,
					ServerMessage: fmt.Sprintf("while scanning a row for the CreateEmailChange function -> %v", err),
					ClientMessage: "An unexpected error has occurred while processing your request.",
					Status:        http.StatusInternalServerError,
				}
			}

			expires, p := parseTime(raw)
			if p != nil {
				return p
			}

			if expires.Before(time.Now()) {
				p = db.DeleteEmailChangeById(id)
				if p != nil {
					return p
				}
			}
		}
	}

	_, err := db.Connection.Exec(
		"INSERT INTO email_chane (token, new_email, expires, user_id) VALUES (?, ?, ?, ?)",
		e.Token, e.NewEmail, e.Expires, e.UserId,
	)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return &problems.Problem{
				Type:          problems.DatabaseProblem,
				ServerMessage: fmt.Sprintf("duplicate entry error while creating an email change token: %v", err),
				ClientMessage: "An email change request already exists for this user.",
				Status:        http.StatusConflict,
			}
		}

		return &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("error while creating a email change token: %v", err),
			ClientMessage: "There was an error while trying to initialize an email change procedure.",
			Status:        http.StatusInternalServerError,
		}
	}

	return nil
}

func (db *Database) GetEmailChangeByToken(token string) (*models.EmailChange, *problems.Problem) {
	var rawTime string
	e := new(models.EmailChange)

	err := db.Connection.QueryRow(
		"SELECT id, token, new_email, expires, user_id FROM email_change WHERE token=?",
		token,
	).Scan(&e.Id, &e.Token, &e.NewEmail, &rawTime, &e.UserId)

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

	expires, p := parseTime(rawTime)
	if p != nil {
		return nil, p
	}

	e.Expires = expires
	return e, nil
}

func (db *Database) DeleteEmailChangeById(id string) *problems.Problem {
	_, err := db.Connection.Exec("DELETE FROM email_change WHERE id=?;", id)
	if err != nil {
		return &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("error while deleting an email change by id: %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
	}

	return nil
}

func (db *Database) GetExpiredEmailChanges() ([]models.EmailChange, *problems.Problem) {
	rows, err := db.Connection.Query("SELECT * FROM email_change WHERE expires <= NOW();")

	if err != nil {
		return nil, &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("error while trying to retrieve expired password changes: %v", err),
			Status:        http.StatusInternalServerError,
		}
	}

	var expired []models.EmailChange
	for rows.Next() {
		var rawTime string
		var emailChange models.EmailChange
		if err := rows.Scan(&emailChange.Id, &emailChange.Token, &emailChange.NewEmail, &rawTime, &emailChange.UserId); err != nil {
			return nil, &problems.Problem{
				Type:          problems.DatabaseProblem,
				ServerMessage: fmt.Sprintf("error while scanning expired email changes: %v", err),
				Status:        http.StatusInternalServerError,
			}
		}

		parsed, err := parseTime(rawTime)
		if err != nil {
			return nil, err
		}

		emailChange.Expires = parsed
		expired = append(expired, emailChange)
	}

	return expired, nil
}
