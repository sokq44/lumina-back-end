package handlers

import (
	"backend/models"
	"backend/utils/problems"
	"encoding/json"
	"fmt"
	"net/http"
)

func CreateArticleComment(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		Comment   models.Comment `json:"comment"`
		ArticleId string         `json:"article_id"`
	}

	var body RequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		p := problems.Problem{
			Type:          problems.HandlerProblem,
			ServerMessage: fmt.Sprintf("(comments endpoint) while decoding the request body -> %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusBadRequest,
		}
		p.Handle(w, r)
		return
	}

	commentId, p := db.CreateComment(body.Comment)
	if p.Handle(w, r) {
		return
	}

	fmt.Println(commentId)

	_, p = db.CreateArticlesComment(body.ArticleId, commentId)
	if p.Handle(w, r) {
		return
	}
}

func UpdateArticleComment(w http.ResponseWriter, r *http.Request) {
	var comment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		p := problems.Problem{
			Type:          problems.HandlerProblem,
			ServerMessage: fmt.Sprintf("(comments endpoint) while decoding the request body -> %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusBadRequest,
		}
		p.Handle(w, r)
		return
	}

	if db.UpdateComment(comment).Handle(w, r) {
		return
	}
}

func DeleteArticleComment(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		Id string `json:"id"`
	}

	var body RequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		p := problems.Problem{
			Type:          problems.HandlerProblem,
			ServerMessage: fmt.Sprintf("(comments endpoint) while decoding the request body -> %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusBadRequest,
		}
		p.Handle(w, r)
		return
	}

	if db.DeleteCommentById(body.Id).Handle(w, r) {
		return
	}
}
