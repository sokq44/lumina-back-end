package database

import (
	"backend/models"
	"backend/utils/problems"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
)

func (db *Database) CreateUser(u models.User) *problems.Problem {
	_, err := db.Connection.Exec(
		"INSERT INTO users (username, email, image_url, password) values (?, ?, ?, ?);",
		u.Username, u.Email, u.ImageUrl, u.Password,
	)

	if err != nil {
		return &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("while creating a new user -> %v", err),
			ClientMessage: "An error occurred while creating a new user.",
			Status:        http.StatusInternalServerError,
		}
	}

	return nil
}

func (db *Database) UpdateUser(u models.User) *problems.Problem {
	_, err := db.Connection.Exec(
		"UPDATE users SET username=?, email=?, image_url=?, password=?, verified=? WHERE id=?",
		u.Username, u.Email, u.ImageUrl, u.Password, u.Verified, u.Id,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("while updating a user -> %v", err),
			ClientMessage: "Error while trying to get user's data.",
			Status:        http.StatusNotFound,
		}
	} else if err != nil {
		return &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("while updating a user -> %v", err),
			ClientMessage: "An error occurred while modifying a user.",
			Status:        http.StatusInternalServerError,
		}
	}

	return nil
}

func (db *Database) DeleteUserById(id string) *problems.Problem {
	_, err := db.Connection.Exec("DELETE FROM users WHERE id=?;", id)

	if err != nil {
		return &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("while deleting a user by id -> %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
	}

	return nil
}

func (db *Database) GetUserById(id string) (*models.User, *problems.Problem) {
	user := &models.User{Id: id}

	err := db.Connection.QueryRow(
		"SELECT username, email, image_url, password, verified FROM users WHERE id=?;",
		id,
	).Scan(&user.Username, &user.Email, &user.ImageUrl, &user.Password, &user.Verified)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("error while getting a user by id: %v", err),
			ClientMessage: "Error while trying to get user's data.",
			Status:        http.StatusNotFound,
		}
	} else if err != nil {
		return nil, &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("error while getting a user by id: %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
	}

	return user, nil
}

func (db *Database) GetUserByEmail(email string) (*models.User, *problems.Problem) {
	user := &models.User{Email: email}

	err := db.Connection.QueryRow(
		"SELECT id, username, image_url, password, verified FROM users WHERE email=?;",
		email,
	).Scan(&user.Id, &user.Username, &user.ImageUrl, &user.Password, &user.Verified)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("error while getting a user by id: %v", err),
			ClientMessage: "Error while trying to get user's data.",
			Status:        http.StatusNotFound,
		}
	} else if err != nil {
		return nil, &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("error while getting a user by email: %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
	}

	return user, nil
}

func (db *Database) UserExists(u models.User) (bool, *problems.Problem) {
	var id string
	err := db.Connection.QueryRow(
		"SELECT id FROM users WHERE username=? or email=?;",
		u.Username, u.Email,
	).Scan(&id)

	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	} else if err != nil {
		return false, &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("error while checking whether a user exists: %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
	}

	return true, nil
}

func (db *Database) VerifyUser(id string) *problems.Problem {
	_, err := db.Connection.Exec(
		"UPDATE users SET verified=TRUE WHERE id=?;",
		id,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("error while verifying a user: %v", err),
			ClientMessage: "Error while trying to get user's data.",
			Status:        http.StatusNotFound,
		}
	} else if err != nil {
		return &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("error while verifying a user: %v", err),
			ClientMessage: "An error has occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
	}

	return nil
}
