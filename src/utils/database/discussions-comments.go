package database

import (
	"backend/models"
	"backend/utils/problems"
	"fmt"
	"net/http"
)

func (db *Database) CreateDiscussionsComment(discussionId string, commentId string) *problems.Problem {
	_, err := db.Connection.Exec(
		"INSERT INTO discussions_comments (discussion_id, comment_id) VALUES (?, ?)",
		discussionId, commentId,
	)
	if err != nil {
		return &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("while trying to create discussion's comment -> %v", err),
			ClientMessage: "An unexpected error has occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
	}

	return nil
}

func (db *Database) GetFirstCommentByDiscussionId(discussionId string) (*models.Comment, *problems.Problem) {
	var commentId string
	err := db.Connection.QueryRow(
		"SELECT comment_id FROM discussions_comments WHERE discussion_id=?",
		discussionId,
	).Scan(&commentId)
	if err != nil {
		return nil, &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("while trying to retrieve comment_id by discussion_id -> %v", err),
			ClientMessage: "An unexpected Error has occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
	}

	return db.GetCommentById(commentId)
}
