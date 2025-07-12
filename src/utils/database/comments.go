package database

import (
	"backend/models"
	"backend/utils/problems"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
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
	var rawCreatedAt string
	var rawLastModified string

	err := db.Connection.QueryRow(
		"SELECT user_id, content, created_at, last_modified FROM comments WHERE id LIKE ?;",
		id,
	).Scan(&comment.UserId, &comment.Content, &rawCreatedAt, &rawLastModified)
	if err != nil {
		return nil, &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("while retrieving a comment by its id -> %v", err),
			ClientMessage: "An unexpected error occurred while retrieving a comment.",
			Status:        http.StatusInternalServerError,
		}
	}

	parsedCreatedAt, p := parseTime(rawCreatedAt)
	if p != nil {
		return nil, p
	}

	parsedLastModified, p := parseTime(rawLastModified)
	if p != nil {
		return nil, p
	}

	comment.CreatedAt = parsedCreatedAt
	comment.LastModified = parsedLastModified

	return comment, nil
}

func (db *Database) GetCommentsByArticleId(id string) ([]models.Comment, *problems.Problem) {
	rows, err := db.Connection.Query(`
	SELECT
		comments.id,
		comments.user_id,
		comments.content,
		comments.created_at,
		comments.last_modified
	FROM
		comments
		JOIN articles_comments ON comments.id = articles_comments.comment_id
	WHERE
		articles_comments.article_id LIKE ? ORDER BY created_at DESC;`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("while fetching comments by article's ID -> %v", err),
			ClientMessage: "There are no comments for this article.",
			Status:        http.StatusNotFound,
		}
	} else if err != nil {
		return nil, &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("while fetching comments by article's ID -> %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
	}

	comments := make([]models.Comment, 0)
	for rows.Next() {
		var rawCreatedAt string
		var rawLastModified string
		var comment models.Comment

		err = rows.Scan(&comment.Id, &comment.UserId, &comment.Content, &rawCreatedAt, &rawLastModified)
		if err != nil {
			return nil, &problems.Problem{
				Type:          problems.DatabaseProblem,
				ServerMessage: fmt.Sprintf("while fetching comments by article's ID and scanning a row -> %v", err),
				ClientMessage: "An error occurred while processing your request.",
				Status:        http.StatusInternalServerError,
			}
		}

		createdAt, p := parseTime(rawCreatedAt)
		if p != nil {
			return nil, p
		}

		lastModified, p := parseTime(rawLastModified)
		if p != nil {
			return nil, p
		}

		comment.CreatedAt = createdAt
		comment.LastModified = lastModified
		comments = append(comments, comment)
	}

	return comments, nil
}

func (db *Database) UpdateComment(c models.Comment) *problems.Problem {
	_, err := db.Connection.Exec(
		"UPDATE comments SET user_id=?, content=?, created_at=?, last_modified=? WHERE id LIKE ?;",
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
