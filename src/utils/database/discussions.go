package database

import (
	"backend/models"
	"backend/utils/problems"
	"fmt"
	"net/http"
	"strconv"
)

func (db *Database) CreateDiscussion(d models.Discussion) (string, *problems.Problem) {
	result, err := db.Connection.Exec("INSERT INTO discussions (created_at) values (?)", d.CreatedAt)
	if err != nil {
		return "", &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("while creating a new discussion -> %v", err),
			ClientMessage: "An unexpected error occurred while creating a new discussion.",
			Status:        http.StatusInternalServerError,
		}
	}

	id, err := result.LastInsertId()
	if err != nil {
		return "", &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("while retrieving a newly created discussion's id -> %v", err),
			ClientMessage: "An unexpected error occurred while creating a new discussion.",
			Status:        http.StatusInternalServerError,
		}
	}

	return strconv.FormatInt(id, 10), nil
}

func (db *Database) GetDiscussionById(id string) (*models.Discussion, *problems.Problem) {
	discussion := &models.Discussion{Id: id}
	var rawTime string

	err := db.Connection.QueryRow("SELECT created_at FROM discussions WHERE id LIKE '?'", id).Scan(rawTime)
	if err != nil {
		return nil, &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("while retrieving a comment by its id -> %v", err),
			ClientMessage: "An unexpected error occurred while retrieving a comment.",
			Status:        http.StatusInternalServerError,
		}
	}

	parsedTime, p := parseTime(rawTime)
	if p == nil {
		discussion.CreatedAt = parsedTime
	} else {
		return nil, p
	}

	return discussion, nil
}

func (db *Database) UpdateDiscussion(d models.Discussion) *problems.Problem {
	_, err := db.Connection.Exec("UPDATE discussions SET created_at=? WHERE id LIKE '?'", d.CreatedAt, d.Id)
	if err != nil {
		return &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("while updating a discussion -> %v", err),
			ClientMessage: "An unexpected error occurred while updating a discussion.",
			Status:        http.StatusInternalServerError,
		}
	}

	return nil
}

func (db *Database) DeleteDiscussion(d models.Discussion) *problems.Problem {
	_, err := db.Connection.Exec("DELETE FROM discussions WHERE id LIKE '?'", d.Id)
	if err != nil {
		return &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("while deleting a discussion -> %v", err),
			ClientMessage: "An unexpected error occurred while deleting a discussion.",
			Status:        http.StatusInternalServerError,
		}
	}

	return nil
}
