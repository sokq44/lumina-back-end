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

	/* Assigning correct user id to the comment structure. */
	_, access, p := jwt.GetRefAccFromRequest(r)
	if p.Handle(w, r) {
		return
	}
	claims, p := jwt.DecodePayload(access)
	if p.Handle(w, r) {
		return
	}
	body.Comment.UserId = claims["user"].(string)

	now := time.Now()
	body.Comment.CreatedAt = now
	body.Comment.LastModified = now
	commentId, p := db.CreateComment(body.Comment)
	if p.Handle(w, r) {
		return
	}

	_, p = db.CreateArticlesComment(body.ArticleId, commentId)
	if p.Handle(w, r) {
		return
	}

	_, err := w.Write([]byte(commentId))
	if err != nil {
		p = &problems.Problem{
			Type:          problems.HandlerProblem,
			ServerMessage: fmt.Sprintf("while trying to write newly created comment's id to the response writer -> %v", err),
			ClientMessage: "An unexpected error has occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
		p.Handle(w, r)
		return
	}
}

func GetAllArticleComments(w http.ResponseWriter, r *http.Request) {
	type UserResponse struct {
		Id       string `json:"id"`
		Email    string `json:"email"`
		ImageUrl string `json:"image"`
		Username string `json:"username"`
	}

	type ResponseElement struct {
		Id           string       `json:"id"`
		User         UserResponse `json:"user"`
		Content      string       `json:"content"`
		CreatedAt    time.Time    `json:"created_at"`
		LastModified time.Time    `json:"last_modified"`
	}

	query := r.URL.Query()
	articleId := query.Get("articleId")

	if articleId == "" {
		if err := json.NewEncoder(w).Encode([]ResponseElement{}); err != nil {
			p := problems.Problem{
				Type:          problems.HandlerProblem,
				ServerMessage: fmt.Sprintf("while trying to encode empty array (endpoint for getting all comments for an article) -> %v", err),
				ClientMessage: "An error occurred while processing your request.",
				Status:        http.StatusInternalServerError,
			}
			p.Handle(w, r)
			return
		}
	}

	comments, p := db.GetCommentsByArticleId(articleId)
	if p.Handle(w, r) {
		return
	}

	var response []ResponseElement
	for _, comment := range comments {
		element := ResponseElement{
			Id:           comment.Id,
			Content:      comment.Content,
			CreatedAt:    comment.CreatedAt,
			LastModified: comment.LastModified,
		}

		user, p := db.GetUserById(comment.UserId)
		if p.Handle(w, r) {
			return
		}

		element.User.Id = user.Id
		element.User.Email = user.Email
		element.User.ImageUrl = user.ImageUrl
		element.User.Username = user.Username

		response = append(response, element)
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		p := problems.Problem{
			Type:          problems.HandlerProblem,
			ServerMessage: fmt.Sprintf("while trying to encode comments (endpoint for getting all comments for an article) -> %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
		p.Handle(w, r)
		return
	}
}

func UpdateArticleComment(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		Content   string `json:"content"`
		CommentId string `json:"comment_id"`
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

	/* Checking whether request sender is the author of the comment. */
	_, access, p := jwt.GetRefAccFromRequest(r)
	if p.Handle(w, r) {
		return
	}

	claims, p := jwt.DecodePayload(access)
	if p.Handle(w, r) {
		return
	}

	comment, p := db.GetCommentById(body.CommentId)
	if p.Handle(w, r) {
		return
	}

	if comment.UserId != claims["user"].(string) {
		p = &problems.Problem{
			Type:          problems.HandlerProblem,
			ServerMessage: fmt.Sprintf("Request sender isn't the author of comment (%s).", comment.Id),
			ClientMessage: "You can't modify this comment.",
			Status:        http.StatusUnauthorized,
		}
		p.Handle(w, r)
	} else {
		comment.Content = body.Content
		comment.LastModified = time.Now()
		db.UpdateComment(*comment).Handle(w, r)
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

	/* Checking whether request sender is the author of the comment. */
	_, access, p := jwt.GetRefAccFromRequest(r)
	if p.Handle(w, r) {
		return
	}

	claims, p := jwt.DecodePayload(access)
	if p.Handle(w, r) {
		return
	}

	commentFromDb, p := db.GetCommentById(body.Id)
	if p.Handle(w, r) {
		return
	}

	if commentFromDb.UserId != claims["user"].(string) {
		p = &problems.Problem{
			Type:          problems.HandlerProblem,
			ServerMessage: fmt.Sprintf("Request sender isn't the author of comment (%s).", body.Id),
			ClientMessage: "You can't delete this comment.",
			Status:        http.StatusUnauthorized,
		}
		p.Handle(w, r)
	} else {
		db.DeleteCommentById(body.Id).Handle(w, r)
	}
}
