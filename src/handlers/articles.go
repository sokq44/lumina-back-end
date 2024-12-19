package handlers

import (
	"backend/models"
	"backend/utils/jwt"
	"backend/utils/problems"
	"encoding/json"
	"net/http"
)

func AddArticle(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	var body RequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		p := problems.Problem{
			Type:          problems.HandlerError,
			ServerMessage: "while decoding the request body -> " + err.Error(),
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

	json.NewEncoder(w).Encode(articles)
}
