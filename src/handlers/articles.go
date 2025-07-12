package handlers

import (
	"backend/config"
	"backend/models"
	"backend/utils/jwt"
	"backend/utils/problems"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func SaveArticle(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		Id        string `json:"id"`
		Title     string `json:"title"`
		Content   string `json:"content"`
		BannerUrl string `json:"banner"`
		Public    bool   `json:"public"`
	}

	var body RequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		p := problems.Problem{
			Type:          problems.HandlerProblem,
			ServerMessage: fmt.Sprintf("while decoding the request body -> %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusBadRequest,
		}
		p.Handle(w, r)
		return
	}

	_, access, p := jwt.GetRefAccFromRequest(r)
	if p.Handle(w, r) {
		return
	}

	claims, p := jwt.DecodePayload(access)
	if p.Handle(w, r) {
		return
	}

	var id string
	if body.Id == "" {
		article := &models.Article{
			Title:     body.Title,
			Public:    body.Public,
			Content:   body.Content,
			BannerUrl: config.Host + "/images/default-banner.png",
			UserId:    claims["user"].(string),
		}
		id, p = db.CreateArticle(article)
		if p.Handle(w, r) {
			return
		}
	} else {
		article := &models.Article{
			Id:        body.Id,
			Title:     body.Title,
			Public:    body.Public,
			Content:   body.Content,
			BannerUrl: body.BannerUrl,
			UserId:    claims["user"].(string),
		}
		if db.UpdateArticle(article).Handle(w, r) {
			return
		}
		id = body.Id
	}

	retrievedArticle, p := db.GetArticleById(id)
	if p.Handle(w, r) {
		return
	}

	if err := json.NewEncoder(w).Encode(retrievedArticle.Id); err != nil {
		p := problems.Problem{
			Type:          problems.HandlerProblem,
			ServerMessage: fmt.Sprintf("while encoding the response body -> %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
		p.Handle(w, r)
		return
	}
}

func GetArticles(w http.ResponseWriter, r *http.Request) {
	type ResponseData struct {
		Id        string    `json:"id"`
		User      string    `json:"user"`
		UserImage string    `json:"user_image"`
		Title     string    `json:"title"`
		Public    bool      `json:"public"`
		BannerUrl string    `json:"banner"`
		Content   string    `json:"content"`
		CreatedAt time.Time `json:"created_at"`
	}

	query := r.URL.Query()
	q := query.Get("q")
	limitstr := query.Get("limit")

	limit, err := strconv.Atoi(limitstr)
	if err != nil {
		p := problems.Problem{
			Type:          problems.HandlerProblem,
			ServerMessage: fmt.Sprintf("Failed to parse limit from string to int -> %v", err),
			ClientMessage: "An unexpected error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
		p.Handle(w, r)
		return
	}

	_, access, p := jwt.GetRefAccFromRequest(r)
	if p.Handle(w, r) {
		return
	}

	claims, p := jwt.DecodePayload(access)
	if p.Handle(w, r) {
		return
	}

	articles, p := db.GetArticlesByUserId(claims["user"].(string), q, limit)
	if p.Handle(w, r) {
		return
	}

	var articlesResponse []ResponseData
	for _, article := range articles {
		user, p := db.GetUserById(article.UserId)
		if p.Handle(w, r) {
			return
		}

		articlesResponse = append(articlesResponse, ResponseData{
			Id:        article.Id,
			Title:     article.Title,
			Public:    article.Public,
			Content:   article.Content,
			BannerUrl: article.BannerUrl,
			CreatedAt: article.CreatedAt,
			User:      user.Username,
			UserImage: user.ImageUrl,
		})
	}

	if err := json.NewEncoder(w).Encode(articlesResponse); err != nil {
		p := problems.Problem{
			Type:          problems.HandlerProblem,
			ServerMessage: fmt.Sprintf("while encoding the response body -> %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
		p.Handle(w, r)
		return
	}
}

func GetArticle(w http.ResponseWriter, r *http.Request) {
	type ResponseData struct {
		Id        string    `json:"id"`
		User      string    `json:"user"`
		UserImage string    `json:"user_image"`
		Title     string    `json:"title"`
		BannerUrl string    `json:"banner"`
		Content   string    `json:"content"`
		Public    bool      `json:"public"`
		CreatedAt time.Time `json:"created_at"`
	}

	query := r.URL.Query()
	id := query.Get("article")

	article, p := db.GetArticleById(id)
	if p.Handle(w, r) {
		return
	}

	user, p := db.GetUserById(article.UserId)
	if p.Handle(w, r) {
		return
	}

	response := ResponseData{
		Id:        article.Id,
		Title:     article.Title,
		Public:    article.Public,
		Content:   article.Content,
		BannerUrl: article.BannerUrl,
		CreatedAt: article.CreatedAt,
		User:      user.Username,
		UserImage: user.ImageUrl,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		p := problems.Problem{
			Type:          problems.HandlerProblem,
			ServerMessage: fmt.Sprintf("while encoding the response body -> %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
		p.Handle(w, r)
		return
	}
}

func GetSuggestedArticles(w http.ResponseWriter, r *http.Request) {
	type ResponseData struct {
		Id        string    `json:"id"`
		User      string    `json:"user"`
		UserImage string    `json:"user_image"`
		Title     string    `json:"title"`
		BannerUrl string    `json:"banner"`
		Content   string    `json:"content"`
		CreatedAt time.Time `json:"created_at"`
	}

	query := r.URL.Query()
	q := query.Get("q")
	limitstr := query.Get("limit")

	limit, err := strconv.Atoi(limitstr)
	if err != nil {
		p := problems.Problem{
			Type:          problems.HandlerProblem,
			ServerMessage: fmt.Sprintf("Failed to parse limit from string to int -> %v", err),
			ClientMessage: "An unexpected error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
		p.Handle(w, r)
		return
	}

	articles, p := db.GetPublicArticles(q, limit)
	if p.Handle(w, r) {
		return
	}

	response := make([]ResponseData, 0)
	for _, article := range articles {
		user, p := db.GetUserByArticleId(article.Id)
		if p.Handle(w, r) {
			return
		}

		response = append(response, ResponseData{
			Id:        article.Id,
			Title:     article.Title,
			Content:   article.Content,
			BannerUrl: article.BannerUrl,
			CreatedAt: article.CreatedAt,
			User:      user.Username,
			UserImage: user.ImageUrl,
		})
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		p = &problems.Problem{
			Type:          problems.HandlerProblem,
			ServerMessage: fmt.Sprintf("error while encoding response for GetSuggestedArticles handler -> %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
		p.Handle(w, r)
		return
	}
}

func DeleteArticle(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		Id string `json:"id"`
	}

	var body RequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		p := problems.Problem{
			Type:          problems.HandlerProblem,
			ServerMessage: fmt.Sprintf("while decoding the request body -> %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusBadRequest,
		}
		p.Handle(w, r)
		return
	}

	if db.DeleteArticlesCommentByArticleId(body.Id).Handle(w, r) {
		return
	}

	if db.DeleteArticleById(body.Id).Handle(w, r) {
		return
	}
}
