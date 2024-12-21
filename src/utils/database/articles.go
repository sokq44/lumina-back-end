package database

import (
	"backend/models"
	"backend/utils/problems"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
)

func (db *Database) CreateArticle(article models.Article) *problems.Problem {
	_, err := db.Connection.Exec("INSERT INTO articles (title, content, user_id) VALUES (?, ?, ?);",
		article.Title, article.Content, article.UserId,
	)
	if err != nil {
		return &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("while creating a new article -> %v", err),
			ClientMessage: "An error occurred while creating a new article.",
			Status:        http.StatusInternalServerError,
		}
	}

	return nil
}

func (db *Database) GetArticlesByUserId(userId string) ([]models.Article, *problems.Problem) {
	rows, err := db.Connection.Query("SELECT id, title, content, created_at FROM articles WHERE user_id = ?;", userId)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("while getting articles by user id -> %v", err),
			ClientMessage: "There are no articles associated with that person.",
			Status:        http.StatusNotFound,
		}
	} else if err != nil {
		return nil, &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("while getting articles by user id -> %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
	}

	articles := make([]models.Article, 0)
	for rows.Next() {
		var article models.Article
		var rawTime string
		if err := rows.Scan(&article.Id, &article.Title, &article.Content, &rawTime); err != nil {
			return nil, &problems.Problem{
				Type:          problems.DatabaseProblem,
				ServerMessage: fmt.Sprintf("while scanning articles -> %v", err),
				ClientMessage: "An error occurred while processing your request.",
				Status:        http.StatusInternalServerError,
			}
		}

		time, p := parseTime(rawTime)
		if p != nil {
			return nil, p
		}

		article.CreatedAt = time
		articles = append(articles, article)
	}

	return articles, nil
}

func (db *Database) DeleteArticleById(id string) *problems.Problem {
	_, err := db.Connection.Exec("DELETE FROM articles WHERE id = ?;", id)
	if errors.Is(err, sql.ErrNoRows) {
		return &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("while deleting article by id -> %v", err),
			ClientMessage: "There is no such article.",
			Status:        http.StatusNotFound,
		}
	} else if err != nil {
		return &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("while deleting article by id -> %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
	}

	return nil
}
