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

func (db *Database) GetArticlesByUserId(userId, phrase string, limit int) ([]models.Article, *problems.Problem) {
	query := `
    SELECT a.id, a.title, a.content, a.created_at, a.public, a.banner_url
    FROM articles a
    WHERE a.user_id = ? AND a.title LIKE ?
    ORDER BY a.created_at DESC
    LIMIT ?;`

	rows, err := db.Connection.Query(query, userId, "%"+phrase+"%", limit)
	if errors.Is(err, sql.ErrNoRows) {
		rows.Close()
		return nil, &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("while getting articles by user id -> %v", err),
			ClientMessage: "There are no articles associated with you.",
			Status:        http.StatusNotFound,
		}
	} else if err != nil {
		rows.Close()
		return nil, &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("while getting articles by user id -> %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
	}
	defer rows.Close()

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

func (db *Database) GetPublicArticles(phrase string, limit int) ([]models.Article, *problems.Problem) {
	query := `
    SELECT a.id, a.title, a.content, a.user_id, a.banner_url, a.created_at
    FROM articles a
    JOIN users u ON a.user_id = u.id
    WHERE a.public=TRUE AND (a.title LIKE ? OR u.username LIKE ?)
    ORDER BY a.created_at DESC
    LIMIT ?;`

	rows, err := db.Connection.Query(query, "%"+phrase+"%", "%"+phrase+"%", limit)
	if errors.Is(err, sql.ErrNoRows) {
		rows.Close()
		return nil, &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("error while retrieving all articles -> %v", err),
			ClientMessage: "No articles have been found.",
			Status:        http.StatusNotFound,
		}
	} else if err != nil {
		rows.Close()
		return nil, &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("error while retrieving all articles -> %v", err),
			ClientMessage: "An error occurred while retrieving articles.",
			Status:        http.StatusInternalServerError,
		}
	}
	defer rows.Close()

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

func (db *Database) GetArticleReadsByArticleId(id string) (int, *problems.Problem) {
	var reads int
	q := "SELECT COUNT(*) FROM articles_reads WHERE article_id=?;"
	err := db.Connection.QueryRow(q, id).Scan(&reads)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	} else if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return -1, &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("while trying to get the number of article's reads -> %v", err),
			ClientMessage: "An error occurred while processing your request.",
		}
	}

	return reads, nil
}

func (db *Database) GetArticleRatingsByArticleId(id string) ([]int, *problems.Problem) {
	q := "SELECT rating FROM articles_ratings WHERE article_id=?;"
	rows, err := db.Connection.Query(q, id)
	if errors.Is(err, sql.ErrNoRows) {
		return []int{}, nil
	} else if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("while trying to get ratings of an article -> %v", err),
			ClientMessage: "An error occurred while processing your request.",
		}
	}
	defer rows.Close()

	ratings := make([]int, 0)
	for rows.Next() {
		var rating int
		if err := rows.Scan(&rating); err != nil {
			return nil, &problems.Problem{
				Type:          problems.DatabaseProblem,
				ServerMessage: fmt.Sprintf("while trying to scan of the ratings of an article -> %v", err),
				ClientMessage: "An error occurred while processing your request.",
			}
		}
		ratings = append(ratings, rating)
	}

	if rows.Err() != nil {
		return nil, &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("error iterating article's rating rows: %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
	}

	return ratings, nil
}

func (db *Database) UpdateArticleReads(articleId, userId string) *problems.Problem {
	var dummy int
	q := "SELECT 1 FROM articles_reads WHERE article_id=? AND user_id=?;"
	err := db.Connection.QueryRow(q, articleId, userId).Scan(&dummy)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("while trying to establish whether a user has already seen an article -> %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
	}

	if errors.Is(err, sql.ErrNoRows) {
		q = "INSERT INTO articles_reads (article_id, user_id) VALUES (? , ?);"
		_, err = db.Connection.Exec(q, articleId, userId)
		if err != nil {
			return &problems.Problem{
				Type:          problems.DatabaseProblem,
				ServerMessage: fmt.Sprintf("while trying to insert a new articles_reads row -> %v", err),
				ClientMessage: "An error occurred while processing your request.",
				Status:        http.StatusInternalServerError,
			}
		}
	}

	return nil
}

func (db *Database) UpdateArticleRatings(articleId, userId string, rating int) *problems.Problem {
	var dummy int
	q := "SELECT 1 FROM articles_ratings WHERE article_id=? AND user_id=?;"
	err := db.Connection.QueryRow(q, articleId, userId).Scan(&dummy)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("while trying to establish whether a user has already seen an article -> %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
	}

	if errors.Is(err, sql.ErrNoRows) {
		q = "INSERT INTO articles_ratings (rating, article_id, user_id) VALUES (? , ?, ?);"
	} else {
		q = "UPDATE articles_ratings SET rating=? WHERE article_id=? AND user_id=?"
	}

	_, err = db.Connection.Exec(q, rating, articleId, userId)
	if err != nil {
		return &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("while trying to insert a new articles_reads row -> %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
	}

	return nil
}
