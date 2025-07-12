package database

import (
	"backend/models"
	"backend/utils/problems"
	"fmt"
	"net/http"

	"github.com/google/uuid"
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

func (db *Database) GetArticleByCommentId(commentId string) (*models.Article, *problems.Problem) {
	var article *models.Article
	var rawTime string
	err := db.Connection.QueryRow(`
		SELECT id, title, content, created_at, user_id, banner_url, public 
		FROM articles 
		JOIN articles_comments ON articles.id=articles_comments.article_id 
		WHERE articles_comments.comment_id=?`,
		commentId,
	).Scan(&article.Id, &article.Title, &article.Content, &rawTime, &article.UserId, &article.BannerUrl, &article.Public)
	if err != nil {
		return nil, &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("while trying to get article by comment's id -> %v", err),
			ClientMessage: "An unexpected error has occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
	}

	parsedTime, p := parseTime(rawTime)
	if p != nil {
		return nil, p
	}

	article.CreatedAt = parsedTime
	return article, nil
}
