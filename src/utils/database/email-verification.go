package database

import (
	"backend/models"
	"backend/utils/problems"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
)

func (db *Database) CreateEmailVerification(e models.EmailVerification) *problems.Problem {
	_, err := db.Connection.Exec(
		"INSERT INTO email_verification (token, expires, user_id) values (?, ?, ?);",
		e.Token, e.Expires, e.UserId,
	)

	if err != nil {
		return &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("error while creating an email verification: %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
	}

	return nil
}

func (db *Database) GetEmailVerificationByToken(token string) (*models.EmailVerification, *problems.Problem) {
	var id string
	var tk string
	var userId string
	var expires string

	err := db.Connection.QueryRow(
		"SELECT id, token, expires, user_id FROM email_verification WHERE token=?;",
		token,
	).Scan(&id, &tk, &expires, &userId)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("error while getting an email verification by token: %v", err),
			ClientMessage: "The verification link is invalid or has expired.",
			Status:        http.StatusNotFound,
		}
	} else if err != nil {
		return nil, &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("error while getting an email verification by token: %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
	}

	expiresTime, e := parseTime(expires)
	if e != nil {
		return nil, e
	}

	emailVerification := &models.EmailVerification{
		Id:      id,
		Token:   tk,
		UserId:  userId,
		Expires: expiresTime,
	}

	return emailVerification, nil
}

func (db *Database) DeleteEmailVerificationById(id string) *problems.Problem {
	_, err := db.Connection.Exec(
		"DELETE FROM email_verification WHERE id=?;",
		id,
	)

	if err != nil {
		return &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("error while deleting an email verification by id: %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
	}

	return nil
}

func (db *Database) GetExpiredEmailVerifications() ([]models.EmailVerification, *problems.Problem) {
	rows, err := db.Connection.Query("SELECT * FROM email_verification WHERE expires <= NOW();")

	if err != nil {
		return nil, &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("error while trying to retrieve expired email verifications: %v", err),
			Status:        http.StatusInternalServerError,
		}
	}

	var expired []models.EmailVerification
	for rows.Next() {
		var verification models.EmailVerification
		var rawTime string
		if err := rows.Scan(&verification.Id, &verification.Token, &rawTime, &verification.UserId); err != nil {
			return nil, &problems.Problem{
				Type:          problems.DatabaseProblem,
				ServerMessage: fmt.Sprintf("error while scanning expired email verifications: %v", err),
				Status:        http.StatusInternalServerError,
			}
		}

		parsed, err := parseTime(rawTime)
		if err != nil {
			return nil, err
		}

		verification.Expires = parsed
		expired = append(expired, verification)
	}

	return expired, nil
}
