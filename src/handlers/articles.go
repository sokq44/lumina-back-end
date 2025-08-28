package handlers

import (
	"backend/config"
	"backend/models"
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
		TLDR      string `json:"tldr"`
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

	user, p := GetUserFromRequest(r)
	if p.Handle(w, r) {
		return
	}

	var id string
	article := &models.Article{
		Title:   body.Title,
		Public:  body.Public,
		Content: body.Content,
		TLDR:    body.TLDR,
		UserId:  user.Id,
	}
	if body.Id == "" {
		article.BannerUrl = config.Host + "/images/default-banner.png"
		id, p = db.CreateArticle(article)
		if p.Handle(w, r) {
			return
		}
	} else {
		article.Id = body.Id
		article.BannerUrl = body.BannerUrl
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
		TLDR      string    `json:"tldr"`
		CreatedAt time.Time `json:"created_at"`
		Reads     int       `json:"reads"`
		Comments  int       `json:"comments"`
		Ratings   []int     `json:"ratings"`
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

	user, p := GetUserFromRequest(r)
	if p.Handle(w, r) {
		return
	}

	articles, p := db.SearchArticlesByUserId(user.Id, q, limit)
	if p.Handle(w, r) {
		return
	}

	var articlesResponse []ResponseData
	for _, article := range articles {
		user, p := db.GetUserById(article.UserId)
		if p.Handle(w, r) {
			return
		}

		ratings, _, p := db.GetArticleRatingsByArticleId(article.Id)
		if p.Handle(w, r) {
			return
		}

		reads, p := db.GetArticleReadsByArticleId(article.Id)
		if p.Handle(w, r) {
			return
		}

		count, p := db.GetCommentsCountByArticleId(article.Id)
		if p.Handle(w, r) {
			return
		}

		articlesResponse = append(articlesResponse, ResponseData{
			Id:        article.Id,
			Title:     article.Title,
			Public:    article.Public,
			Content:   article.Content,
			TLDR:      article.TLDR,
			BannerUrl: article.BannerUrl,
			CreatedAt: article.CreatedAt,
			User:      user.Username,
			UserImage: user.ImageUrl,
			Ratings:   ratings,
			Reads:     reads,
			Comments:  count,
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
		Reads     int       `json:"reads"`
		Comments  int       `json:"comments"`
		MyRating  int       `json:"my_rating"`
		Ratings   []int     `json:"ratings"`
		Id        string    `json:"id"`
		User      string    `json:"user"`
		UserImage string    `json:"user_image"`
		Title     string    `json:"title"`
		BannerUrl string    `json:"banner"`
		Content   string    `json:"content"`
		TLDR      string    `json:"tldr"`
		Public    bool      `json:"public"`
		CreatedAt time.Time `json:"created_at"`
	}

	query := r.URL.Query()
	id := query.Get("article")

	sessionUser, p := GetUserFromRequest(r)
	if p.Handle(w, r) {
		return
	}

	if sessionUser != nil && p == nil {
		if db.UpdateArticleReads(id, sessionUser.Id).Handle(w, r) {
			return
		}
	}

	article, p := db.GetArticleById(id)
	if p.Handle(w, r) {
		return
	}

	user, p := db.GetUserById(article.UserId)
	if p.Handle(w, r) {
		return
	}

	ratings, userIds, p := db.GetArticleRatingsByArticleId(article.Id)
	if p.Handle(w, r) {
		return
	}

	myRating := -1
	if sessionUser != nil {
		for i := 0; i < len(ratings); i++ {
			if userIds[i] == sessionUser.Id {
				myRating = ratings[i]
			}
		}
	}

	reads, p := db.GetArticleReadsByArticleId(article.Id)
	if p.Handle(w, r) {
		return
	}

	count, p := db.GetCommentsCountByArticleId(article.Id)
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
		TLDR:      article.TLDR,
		User:      user.Username,
		UserImage: user.ImageUrl,
		Ratings:   ratings,
		Reads:     reads,
		Comments:  count,
		MyRating:  myRating,
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
		TLDR      string    `json:"tldr"`
		CreatedAt time.Time `json:"created_at"`
		Reads     int       `json:"reads"`
		Comments  int       `json:"comments"`
		Ratings   []int     `json:"ratings"`
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

	articles, p := db.SearchPublicArticles(q, limit)
	if p.Handle(w, r) {
		return
	}

	response := make([]ResponseData, 0)
	for _, article := range articles {
		user, p := db.GetUserByArticleId(article.Id)
		if p.Handle(w, r) {
			return
		}

		ratings, _, p := db.GetArticleRatingsByArticleId(article.Id)
		if p.Handle(w, r) {
			return
		}

		reads, p := db.GetArticleReadsByArticleId(article.Id)
		if p.Handle(w, r) {
			return
		}

		count, p := db.GetCommentsCountByArticleId(article.Id)
		if p.Handle(w, r) {
			return
		}

		response = append(response, ResponseData{
			Id:        article.Id,
			Title:     article.Title,
			Content:   article.Content,
			TLDR:      article.TLDR,
			BannerUrl: article.BannerUrl,
			CreatedAt: article.CreatedAt,
			User:      user.Username,
			UserImage: user.ImageUrl,
			Ratings:   ratings,
			Reads:     reads,
			Comments:  count,
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

func UpdateArticleReads(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		ArticleId string `json:"article_id"`
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

	user, p := GetUserFromRequest(r)
	if p.Handle(w, r) {
		return
	}

	if db.UpdateArticleReads(body.ArticleId, user.Id).Handle(w, r) {
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func UpdateArticleRatings(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		ArticleId string `json:"article_id"`
		Rating    int    `json:"rating"`
	}

	var body RequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		p := problems.Problem{
			Type:          problems.HandlerProblem,
			ServerMessage: "Wrong rating was provided (< 0 or > 5)",
			ClientMessage: "Please provide a rating between 1 and 5. Ratings outside this range are not accepted.",
			Status:        http.StatusBadRequest,
		}
		p.Handle(w, r)
		return
	}

	if body.Rating < 0 || body.Rating > 5 {
		p := problems.Problem{
			Type:          problems.HandlerProblem,
			ServerMessage: "Wrong rating was provided (< 0 or > 5)",
			ClientMessage: "Article may only be rated between 1 and 5.",
			Status:        http.StatusBadRequest,
		}
		p.Handle(w, r)
		return
	}

	user, p := GetUserFromRequest(r)
	if p.Handle(w, r) {
		return
	}

	if db.UpdateArticleRatings(body.ArticleId, user.Id, body.Rating).Handle(w, r) {
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
