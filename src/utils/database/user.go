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
		"INSERT INTO users (id, username, email, image_url, password) values (?, ?, ?, ?, ?);",
		u.Id, u.Username, u.Email, u.ImageUrl, u.Password,
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
	q := "UPDATE users SET username=?, email=?, image_url=?, password=?, verified=?, bio=?, favourites=? WHERE id=?"
	_, err := db.Connection.Exec(q, u.Username, u.Email, u.ImageUrl, u.Password, u.Verified, u.Bio, u.Favourites, u.Id)

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
			ClientMessage: "An unexpected error has occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
	}

	return nil
}

func (db *Database) UserAddSocial(social models.UserSocialPlatform) *problems.Problem {
	if social.Type == models.UnknownST {
		return &problems.Problem{
			Type:          problems.ModelProblem,
			ServerMessage: "invalid social platform type provided",
			ClientMessage: "Invalid social platform type.",
			Status:        http.StatusBadRequest,
		}
	}

	var id string
	q := "SELECT id FROM user_socials WHERE user_id=? AND social_type=?;"
	err := db.Connection.QueryRow(q, social.UserId, social.Type).Scan(&id)
	if err == nil {
		return &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("social platform %s already added for user %s", models.SocialPlatformGetName(social.Type), social.UserId),
			ClientMessage: "This social platform is already linked to your account.",
			Status:        http.StatusConflict,
		}
	} else if !errors.Is(err, sql.ErrNoRows) {
		return &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("error while checking user_socials -> %v", err),
			ClientMessage: "An unexpected error occurred while linking social platform.",
			Status:        http.StatusInternalServerError,
		}
	}

	q = "INSERT INTO user_socials (user_id, social_type, value, label) VALUES (?, ?, ?, ?);"
	_, err = db.Connection.Exec(q, social.UserId, social.Type, social.Value, social.Label)
	if err != nil {
		return &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("while trying to link a social platform to a user -> %v", err),
			ClientMessage: "An unexpected error has occurred when trying to add this social platform to your account.",
			Status:        http.StatusInternalServerError,
		}
	}

	return nil
}

func (db *Database) UserRemoveSocial(userId string, socialType models.SocialType) *problems.Problem {
	q := "DELETE FROM user_socials WHERE user_id=? AND social_type=?;"
	_, err := db.Connection.Exec(q, userId, socialType)
	if err != nil {
		return &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("while trying to remove a social platform link -> %v", err),
			ClientMessage: "An unexpected error has occured when trying to unlink a social platform.",
			Status:        http.StatusInternalServerError,
		}
	}
	return nil
}
