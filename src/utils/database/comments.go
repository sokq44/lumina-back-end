package database

import (
	"backend/models"
	"backend/utils/problems"
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

func (db *Database) CreateComment(c models.Comment) (string, *problems.Problem) {
	id := uuid.New().String()
	_, err := db.Connection.Exec(
		"INSERT INTO comments (id, user_id, content, created_at, last_modified) values (?, ?, ?, ?, ?)",
		id, c.UserId, c.Content, c.CreatedAt, c.LastModified,
	)
	if err != nil {
		return "", &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("while creating a new comment -> %v", err),
			ClientMessage: "An unexpected error occurred while creating a new comment.",
			Status:        http.StatusInternalServerError,
		}
	}

	return id, nil
}

func (db *Database) GetCommentById(id string) (*models.Comment, *problems.Problem) {
	comment := &models.Comment{Id: id}
	var rawTime string

	err := db.Connection.QueryRow(
		"SELECT (user_id, content, created_at, last_modified) FROM comments WHERE id LIKE ?",
		id,
	).Scan(&comment.UserId, &comment.Content, rawTime, &comment.LastModified)
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
		comment.CreatedAt = parsedTime
	} else {
		return nil, p
	}

	return comment, nil
}

func (db *Database) UpdateComment(c models.Comment) *problems.Problem {
	_, err := db.Connection.Exec(
		"UPDATE comments SET user_id=?, content=?, created_at=?, last_modified=? WHERE id LIKE ?",
		c.UserId, c.Content, c.CreatedAt, c.LastModified, c.Id,
	)
	if err != nil {
		return &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("while updating a comment -> %v", err),
			ClientMessage: "An unexpected error occurred while updating a comment.",
			Status:        http.StatusInternalServerError,
		}
	}

	return nil
}

func (db *Database) DeleteCommentById(id string) *problems.Problem {
	_, err := db.Connection.Exec("DELETE FROM comments WHERE id LIKE ?", id)
	if err != nil {
		return &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("while deleting a comment -> %v", err),
			ClientMessage: "An unexpected error occurred while deleting a comment.",
			Status:        http.StatusInternalServerError,
		}
	}

	return nil
}
