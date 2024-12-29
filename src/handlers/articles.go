package handlers

import (
	"backend/models"
	"backend/utils/jwt"
	"backend/utils/problems"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func AddArticle(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		Title   string `json:"title"`
		Content string `json:"content"`
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

	article := models.Article{
		Title:   body.Title,
		Content: body.Content,
		UserId:  claims["user"].(string),
	}
	if db.CreateArticle(article).Handle(w, r) {
		return
	}
}

func GetArticles(w http.ResponseWriter, r *http.Request) {
	type ResponseData struct {
		Id        string    `json:"id"`
		Title     string    `json:"title"`
		Content   string    `json:"content"`
		CreatedAt time.Time `json:"created_at"`
		User      string    `json:"user"`
	}

	_, access, p := jwt.GetRefAccFromRequest(r)
	if p.Handle(w, r) {
		return
	}

	claims, p := jwt.DecodePayload(access)
	if p.Handle(w, r) {
		return
	}

	articles, p := db.GetArticlesByUserId(claims["user"].(string))
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
			Content:   article.Content,
			CreatedAt: article.CreatedAt,
			User:      user.Username,
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
		Title     string    `json:"title"`
		Content   string    `json:"content"`
		CreatedAt time.Time `json:"created_at"`
		User      string    `json:"user"`
	}

	query := r.URL.Query()
	id := query.Get("article")

	log.Println(id)

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
		Content:   article.Content,
		CreatedAt: article.CreatedAt,
		User:      user.Username,
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

	if db.DeleteArticleById(body.Id).Handle(w, r) {
		return
	}
}
