package database

import (
	"backend/models"
	"backend/utils/problems"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

func (db *Database) CreateArticle(article *models.Article) (string, *problems.Problem) {
	id := uuid.New().String()
	_, err := db.Connection.Exec("INSERT INTO articles (id, title, content, user_id, banner_url) VALUES (?, ?, ?, ?, ?);",
		id, article.Title, article.Content, article.UserId, article.BannerUrl,
	)
	if err != nil {
		return "", &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("while creating a new article -> %v", err),
			ClientMessage: "An error occurred while creating a new article.",
			Status:        http.StatusInternalServerError,
		}
	}

	return id, nil
}

func (db *Database) UpdateArticle(article *models.Article) *problems.Problem {
	_, err := db.Connection.Exec("UPDATE articles SET title = ?, content = ?, user_id = ?, public = ?, banner_url = ? WHERE id = ?;",
		article.Title, article.Content, article.UserId, article.Public, article.BannerUrl,
		article.Id,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("while updating article -> %v", err),
			ClientMessage: "There is no such article.",
			Status:        http.StatusNotFound,
		}
	} else if err != nil {
		return &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("while updating article -> %v", err),
			ClientMessage: "An error occurred while updating the article.",
			Status:        http.StatusInternalServerError,
		}
	}

	return nil
}

func (db *Database) GetArticleById(id string) (*models.Article, *problems.Problem) {
	article := new(models.Article)
	var rawTime string

	err := db.Connection.QueryRow(
		"SELECT id, title, content, created_at, user_id, banner_url, public FROM articles WHERE id = ?;",
		id,
	).Scan(&article.Id, &article.Title, &article.Content, &rawTime, &article.UserId, &article.BannerUrl, &article.Public)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("while getting article by id -> %v", err),
			ClientMessage: "There is no such article.",
			Status:        http.StatusNotFound,
		}
	} else if err != nil {
		return nil, &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("while getting article by id -> %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
	}

	time, p := parseTime(rawTime)
	if p != nil {
		return nil, p
	}

	article.CreatedAt = time

	return article, nil
}

func (db *Database) GetArticleByTitle(title string) (*models.Article, *problems.Problem) {
	article := new(models.Article)
	var rawTime string

	err := db.Connection.QueryRow(
		"SELECT id, title, content, created_at, user_id, banner_url, public FROM articles WHERE title = ?;",
		title,
	).Scan(&article.Id, &article.Title, &article.Content, &rawTime, &article.UserId, &article.BannerUrl, &article.Public)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("while getting article by id -> %v", err),
			ClientMessage: "There is no such article.",
			Status:        http.StatusNotFound,
		}
	} else if err != nil {
		return nil, &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("while getting article by id -> %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
	}

	time, p := parseTime(rawTime)
	if p != nil {
		return nil, p
	}

	article.CreatedAt = time

	return article, nil
}

func (db *Database) GetArticlesByUserId(userId string) ([]models.Article, *problems.Problem) {
	rows, err := db.Connection.Query("SELECT id, title, content, created_at, public, banner_url FROM articles WHERE user_id = ?;", userId)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("while getting articles by user id -> %v", err),
			ClientMessage: "There are no articles associated with you.",
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
		if err := rows.Scan(&article.Id, &article.Title, &article.Content, &rawTime, &article.Public, &article.BannerUrl); err != nil {
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
		article.UserId = userId
		articles = append(articles, article)
	}

	return articles, nil
}

func (db *Database) GetPublicArticles() ([]models.Article, *problems.Problem) {
	rows, err := db.Connection.Query("SELECT id, title, content, user_id, banner_url, created_at FROM articles WHERE public=TRUE ORDER BY created_at DESC;")
	if errors.Is(err, sql.ErrNoRows) {
		return nil, &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("error while retrieving all articles -> %v", err),
			ClientMessage: "No articles have been found.",
			Status:        http.StatusNotFound,
		}
	} else if err != nil {
		return nil, &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("error while retrieving all articles -> %v", err),
			ClientMessage: "An error occurred while retrieving articles.",
			Status:        http.StatusInternalServerError,
		}
	}

	articles := make([]models.Article, 0)
	for rows.Next() {
		var article models.Article
		var rawTime string
		if err := rows.Scan(&article.Id, &article.Title, &article.Content, &article.UserId, &article.BannerUrl, &rawTime); err != nil {
			return nil, &problems.Problem{
				Type:          problems.DatabaseProblem,
				ServerMessage: fmt.Sprintf("error while scanning an article row -> %v", err),
				ClientMessage: "An error occurred while retrieving articles.",
				Status:        http.StatusInternalServerError,
			}
		}

		time, p := parseTime(rawTime)
		if p != nil {
			return nil, p
		}

		article.CreatedAt = time
		article.Public = true
		articles = append(articles, article)
	}

	return articles, nil
}

func (db *Database) GetUserByArticleId(id string) (*models.User, *problems.Problem) {
	article, p := db.GetArticleById(id)
	if p != nil {
		return nil, p
	}

	user := new(models.User)
	err := db.Connection.QueryRow(
		"SELECT username, email, password, verified, image_url FROM users WHERE id=?",
		article.UserId,
	).Scan(&user.Username, &user.Email, &user.Password, &user.Verified, &user.ImageUrl)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("error while retrieving user by article id -> %v", err),
			ClientMessage: "Couldn't find any user affiliated with certain article.",
			Status:        http.StatusNotFound,
		}
	} else if err != nil {
		return nil, &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("error while retrieving user by article id -> %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
	}

	return user, nil
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
