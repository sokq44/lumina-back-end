package database

import (
	"backend/models"
	"backend/utils/problems"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
)

func (db *Database) CreateRefreshToken(token models.RefreshToken) *problems.Problem {
	_, err := db.Connection.Exec(
		"INSERT INTO refresh_tokens (id, token, expires, user_id) values(?, ?, ?, ?)",
		token.Id, token.Token, token.Expires, token.UserId,
	)

	if err != nil {
		return &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("error while creating a refresh token: %v", err),
			ClientMessage: "An error occurred while trying to store your session.",
			Status:        http.StatusInternalServerError,
		}
	}

	return nil
}

func (db *Database) GetRefreshTokenByUserId(userId string) (*models.RefreshToken, *problems.Problem) {
	var token models.RefreshToken
	var rawTime string

	err := db.Connection.QueryRow(
		"SELECT * FROM refresh_tokens WHERE user_id=?;",
		userId,
	).Scan(&token.Id, &token.Token, &rawTime, &token.UserId)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("while getting a refresh token by user id: %v", err),
			ClientMessage: "There's no session associated with the provided user.",
			Status:        http.StatusNotFound,
		}
	} else if err != nil {
		return nil, &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("while getting a refresh token by user id: %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
	}

	t, e := parseTime(rawTime)
	if e != nil {
		return nil, e
	}

	token.Expires = t

	return &token, nil
}

func (db *Database) DeleteRefreshTokenById(id string) *problems.Problem {
	_, err := db.Connection.Exec(
		"DELETE FROM refresh_tokens WHERE id=?;",
		id,
	)

	if err != nil {
		return &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("error while deleting a refresh token by id: %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
	}

	return nil
}

func (db *Database) DeleteRefreshTokenByToken(token string) *problems.Problem {
	_, err := db.Connection.Exec(
		"DELETE FROM refresh_tokens WHERE token=?;",
		token,
	)

	if err != nil {
		return &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("error while deleting a refresh token by token: %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
	}

	return nil
}

func (db *Database) GetExpiredRefreshTokens() ([]models.RefreshToken, *problems.Problem) {
	rows, err := db.Connection.Query("SELECT * FROM refresh_tokens WHERE expires <= NOW();")

	if err != nil {
		return nil, &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("error while trying to retrieve expired refresh tokens: %v", err),
			Status:        http.StatusInternalServerError,
		}
	}

	var expired []models.RefreshToken
	for rows.Next() {
		var token models.RefreshToken
		var rawTime string
		if err := rows.Scan(&token.Id, &token.Token, &rawTime, &token.UserId); err != nil {
			return nil, &problems.Problem{
				Type:          problems.DatabaseProblem,
				ServerMessage: fmt.Sprintf("error while scanning expired refresh tokens: %v", err),
				Status:        http.StatusInternalServerError,
			}
		}

		parsed, err := parseTime(rawTime)
		if err != nil {
			return nil, err
		}

		token.Expires = parsed
		expired = append(expired, token)
	}

	return expired, nil
}
