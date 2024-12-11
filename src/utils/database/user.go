package database

import (
	"backend/models"
	"backend/utils/errhandle"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
)

func (db *Database) CreateUser(u models.User) *errhandle.Error {
	_, err := db.Connection.Exec(
		"INSERT INTO users (id, username, email, password) values (?, ?, ?, ?);",
		u.Id, u.Username, u.Email, u.Password,
	)

	if err != nil {
		return &errhandle.Error{
			Type:          errhandle.DatabaseError,
			ServerMessage: fmt.Sprintf("while creating a new user -> %v", err),
			ClientMessage: "An error occurred while creating a new user.",
			Status:        http.StatusInternalServerError,
		}
	}

	return nil
}

func (db *Database) UpdateUser(u models.User) *errhandle.Error {
	_, err := db.Connection.Exec(
		"UPDATE users SET username=?, email=?, password=?, verified=? WHERE id=?",
		u.Username, u.Email, u.Password, u.Verified, u.Id,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return &errhandle.Error{
			Type:          errhandle.DatabaseError,
			ServerMessage: fmt.Sprintf("while updating a user -> %v", err),
			ClientMessage: "Error while trying to get user's data.",
			Status:        http.StatusNotFound,
		}
	} else if err != nil {
		return &errhandle.Error{
			Type:          errhandle.DatabaseError,
			ServerMessage: fmt.Sprintf("while updating a user -> %v", err),
			ClientMessage: "An error occurred while modifying a user.",
			Status:        http.StatusInternalServerError,
		}
	}

	return nil
}

func (db *Database) DeleteUserById(id string) *errhandle.Error {
	_, err := db.Connection.Exec("DELETE FROM users WHERE id=?;", id)

	if err != nil {
		return &errhandle.Error{
			Type:          errhandle.DatabaseError,
			ServerMessage: fmt.Sprintf("while deleting a user by id -> %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
	}

	return nil
}

func (db *Database) GetUserById(id string) (*models.User, *errhandle.Error) {
	user := &models.User{Id: id}

	err := db.Connection.QueryRow(
		"SELECT username, email, verified FROM users WHERE id=?;",
		id,
	).Scan(&user.Username, &user.Email, &user.Verified)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, &errhandle.Error{
			Type:          errhandle.DatabaseError,
			ServerMessage: fmt.Sprintf("error while getting a user by id: %v", err),
			ClientMessage: "Error while trying to get user's data.",
			Status:        http.StatusNotFound,
		}
	} else if err != nil {
		return nil, &errhandle.Error{
			Type:          errhandle.DatabaseError,
			ServerMessage: fmt.Sprintf("error while getting a user by id: %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
	}

	return user, nil
}

func (db *Database) GetUserByEmail(email string) (*models.User, *errhandle.Error) {
	var id string
	var username string
	var password string
	var verified bool

	err := db.Connection.QueryRow(
		"SELECT id, username, password, verified FROM users WHERE email=?;",
		email,
	).Scan(&id, &username, &password, &verified)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, &errhandle.Error{
			Type:          errhandle.DatabaseError,
			ServerMessage: fmt.Sprintf("error while getting a user by id: %v", err),
			ClientMessage: "Error while trying to get user's data.",
			Status:        http.StatusNotFound,
		}
	} else if err != nil {
		return nil, &errhandle.Error{
			Type:          errhandle.DatabaseError,
			ServerMessage: fmt.Sprintf("error while getting a user by email: %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
	}

	user := &models.User{
		Id:       id,
		Username: username,
		Password: password,
		Verified: verified,
	}

	return user, nil
}

func (db *Database) UserExists(u models.User) (bool, *errhandle.Error) {
	var id string

	err := db.Connection.QueryRow(
		"SELECT id FROM users WHERE username=? or email=?;",
		u.Username, u.Email,
	).Scan(&id)

	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	} else if err != nil {
		return false, &errhandle.Error{
			Type:          errhandle.DatabaseError,
			ServerMessage: fmt.Sprintf("error while checking whether a user exists: %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
	}

	return true, nil
}

func (db *Database) VerifyUser(id string) *errhandle.Error {
	_, err := db.Connection.Exec(
		"UPDATE users SET verified=TRUE WHERE id=?;",
		id,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return &errhandle.Error{
			Type:          errhandle.DatabaseError,
			ServerMessage: fmt.Sprintf("error while verifying a user: %v", err),
			ClientMessage: "Error while trying to get user's data.",
			Status:        http.StatusNotFound,
		}
	} else if err != nil {
		return &errhandle.Error{
			Type:          errhandle.DatabaseError,
			ServerMessage: fmt.Sprintf("error while verifying a user: %v", err),
			ClientMessage: "An error has occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
	}

	return nil
}
