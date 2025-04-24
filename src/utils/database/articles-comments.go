package database

import (
	"backend/utils/problems"
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

func (db *Database) CreateArticlesComment(articleId string, commentId string) (string, *problems.Problem) {
	id := uuid.New().String()
	_, err := db.Connection.Exec(
		"INSERT INTO articles_comments (id, article_id, comment_id) VALUES (?, ?, ?)",
		id, articleId, commentId,
	)
	if err != nil {
		return "", &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("while creating a new article's comment -> %v", err),
			ClientMessage: "An unexpected error occurred while commenting on an article.",
			Status:        http.StatusInternalServerError,
		}
	}

	return id, nil
}

func (db *Database) DeleteArticlesCommentByArticleId(articleId string) *problems.Problem {
	_, err := db.Connection.Exec("DELETE FROM articles_comments WHERE article_id LIKE ?", articleId)
	if err != nil {
		return &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("while deleting article's comment row by article_id -> %v", err),
			ClientMessage: "An unexpected error occurred while deleting a comments.",
			Status:        http.StatusInternalServerError,
		}
	}

	return nil
}

func (db *Database) DeleteArticlesCommentByCommentId(commentId string) *problems.Problem {
	_, err := db.Connection.Exec("DELETE FROM articles_comments WHERE comment_id LIKE ?", commentId)
	if err != nil {
		return &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("while deleting article's comment row by comment_id -> %v", err),
			ClientMessage: "An unexpected error occurred while deleting a comments.",
			Status:        http.StatusInternalServerError,
		}
	}

	return nil
}
