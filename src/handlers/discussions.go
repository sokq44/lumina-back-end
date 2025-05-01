package handlers

import (
	"backend/models"
	"backend/utils/jwt"
	"backend/utils/problems"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func CreateArticleDiscussion(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		PrevId  string         `json:"prev_id"`
		Comment models.Comment `json:"comment"`
	}

	type ResponseBody struct {
		CommentId    string `json:"comment_id"`
		DiscussionId string `json:"discussion_id"`
	}

	var body RequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		p := problems.Problem{
			Type:          problems.HandlerProblem,
			ServerMessage: fmt.Sprintf("while trying to read request body -> %v", err),
			ClientMessage: "An unexpected error has occurred while processing your request.",
			Status:        http.StatusBadRequest,
		}
		p.Handle(w, r)
		return
	}

	/* Checking whether a comment with the given id exists. */
	previousComment, p := db.GetCommentById(body.PrevId)
	if p.Handle(w, r) {
		return
	}

	if previousComment == nil {
		p = &problems.Problem{
			Type:          problems.HandlerProblem,
			ServerMessage: "while trying to create a discussion -> previous comment doesn't exist",
			ClientMessage: "The comment which you are trying to answear does not exist anymore.",
			Status:        http.StatusNotFound,
		}
		p.Handle(w, r)
		return
	}

	/* Retrieving the article on which the discussion is taking place. */
	article, p := db.GetArticleByCommentId(body.PrevId)
	if p.Handle(w, r) {
		return
	}

	/* Assigning correct data to the new comment's structure. */
	_, access, p := jwt.GetRefAccFromRequest(r)
	if p.Handle(w, r) {
		return
	}
	claims, p := jwt.DecodePayload(access)
	if p.Handle(w, r) {
		return
	}

	now := time.Now()
	body.Comment.CreatedAt = now
	body.Comment.LastModified = now
	body.Comment.UserId = claims["user"].(string)

	commentId, p := db.CreateComment(body.Comment)
	if p.Handle(w, r) {
		return
	}

	/* Referencing the comment to the appropriate article. */
	_, p = db.CreateArticlesComment(article.Id, commentId)
	if p.Handle(w, r) {
		return
	}

	/* Creating discussion row in the database. */
	discussionId, p := db.CreateDiscussion(models.Discussion{CreatedAt: now})
	if p.Handle(w, r) {
		return
	}

	/* Referencing newly created comment to the newly created discussion */
	if db.CreateDiscussionsComment(discussionId, commentId).Handle(w, r) {
		return
	}

	response := ResponseBody{
		CommentId:    commentId,
		DiscussionId: discussionId,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		p = &problems.Problem{
			Type:          problems.HandlerProblem,
			ServerMessage: fmt.Sprintf("while trying to encode the response after creating a discussion -> %v", err),
			ClientMessage: "An unexpected error has occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
		p.Handle(w, r)
		return
	} else {
		w.WriteHeader(http.StatusCreated)
	}
}

func UpdateArticleDiscussion(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		DiscussionId string         `json:"discussion_id"`
		Comment      models.Comment `json:"comment"`
	}

	var body RequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		p := problems.Problem{
			Type:          problems.HandlerProblem,
			ServerMessage: fmt.Sprintf("while trying to decode request body -> %v", err),
			ClientMessage: "An unexpected error has occurred while processing your request.",
			Status:        http.StatusBadRequest,
		}
		p.Handle(w, r)
		return
	}

	/* Checking whether discussion with the given id exists. */
	discussion, p := db.GetDiscussionById(body.DiscussionId)
	if p.Handle(w, r) {
		return
	}

	if discussion == nil {
		p = &problems.Problem{
			Type:          problems.HandlerProblem,
			ServerMessage: fmt.Sprintf("while trying to update discussion -> discussion with the provided id doesn't exist (%v)", body.DiscussionId),
			ClientMessage: "The discussion in which you are trying to participate doesn't exist anymore.",
			Status:        http.StatusNotFound,
		}
		p.Handle(w, r)
		return
	}

	/* Retrieving first comment from a discussion. */
	firstComment, p := db.GetFirstCommentByDiscussionId(body.DiscussionId)
	if p.Handle(w, r) {
		return
	}

	/* Retrieving article on which the discussion is held. */
	article, p := db.GetArticleByCommentId(firstComment.Id)
	if p.Handle(w, r) {
		return
	}

	/* Assigning correct data to the new comment's structure. */
	_, access, p := jwt.GetRefAccFromRequest(r)
	if p.Handle(w, r) {
		return
	}
	claims, p := jwt.DecodePayload(access)
	if p.Handle(w, r) {
		return
	}

	now := time.Now()
	body.Comment.CreatedAt = now
	body.Comment.LastModified = now
	body.Comment.UserId = claims["user"].(string)

	commentId, p := db.CreateComment(body.Comment)
	if p.Handle(w, r) {
		return
	}

	/* Referencing the comment to the appropriate article. */
	_, p = db.CreateArticlesComment(article.Id, commentId)
	if p.Handle(w, r) {
		return
	}

	/* Referencing newly created comment to the newly created discussion */
	if db.CreateDiscussionsComment(body.DiscussionId, commentId).Handle(w, r) {
		return
	}

	if err := json.NewEncoder(w).Encode(commentId); err != nil {
		p = &problems.Problem{
			Type:          problems.HandlerProblem,
			ServerMessage: "while trying to encode newly created comment's id into the http response",
			ClientMessage: "An unexpected error has occurred while trying to process your request.",
			Status:        http.StatusInternalServerError,
		}
		p.Handle(w, r)
		return
	}
}
