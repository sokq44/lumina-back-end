package handlers

import (
	"backend/models"
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
	now := time.Now()
	user, p := GetUserFromRequest(r)
	if p.Handle(w, r) {
		return
	}

	body.Comment.UserId = user.Id
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

	response := make([]map[string]interface{}, 0)
	for _, comment := range comments {
		element := map[string]interface{}{
			"id":            comment.Id,
			"content":       comment.Content,
			"created_at":    comment.CreatedAt,
			"last_modified": comment.LastModified,
			"user":          map[string]string{},
		}

		user, p := db.GetUserById(comment.UserId)
		if p.Handle(w, r) {
			return
		}

		element["user"] = map[string]string{
			"id":       user.Id,
			"email":    user.Email,
			"image":    user.ImageUrl,
			"username": user.Username,
		}

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
	user, p := GetUserFromRequest(r)
	if p.Handle(w, r) {
		return
	}

	comment, p := db.GetCommentById(body.CommentId)
	if p.Handle(w, r) {
		return
	}

	if comment.UserId != user.Id {
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
	query := r.URL.Query()
	id := query.Get("id")

	/* Checking whether request sender is the author of the comment. */
	user, p := GetUserFromRequest(r)
	if p.Handle(w, r) {
		return
	}

	commentFromDb, p := db.GetCommentById(id)
	if p.Handle(w, r) {
		return
	}

	if commentFromDb.UserId != user.Id {
		p = &problems.Problem{
			Type:          problems.HandlerProblem,
			ServerMessage: fmt.Sprintf("Request sender isn't the author of comment (%s).", id),
			ClientMessage: "You can't delete this comment.",
			Status:        http.StatusUnauthorized,
		}
		p.Handle(w, r)
	} else {
		db.DeleteArticlesCommentByCommentId(id).Handle(w, r)
		db.DeleteCommentById(id).Handle(w, r)
	}
}
