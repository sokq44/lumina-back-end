package handlers

import (
	"backend/models"
	"backend/utils/problems"
	"encoding/json"
	"net/http"
)

func GetProfile(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	username := query.Get("user")

	// TODO: Add hometown location to the user table
	user, p := db.GetUserByUsername(username)
	if p.Handle(w, r) {
		return
	}

	// TODO: Write the functionality
	// articlesCount, p := db.GetArticlesCountByUserId(user.Id)
	// if p.Handle(w, r) {
	// 	return
	// }

	// TODO: Write the functionality
	// wordsCount, p := db.GetWordsCountByUserId(user.Id)
	// if p.Handle(w, r) {
	// 	return
	// }

	// TODO: Write the functionality
	// readsCount, p := db.GetReadsCountByUserId(user.Id)
	// if p.Handle(w, r) {
	// 	return
	// }

	// TODO: Write the functionality
	// avgRating, p := db.GetAvgRatingByUserId(user.Id)
	// if p.Handle(w, r) {
	// 	return
	// }

	// TODO: Write the functionality
	// book, p := db.GetCurrentBookByUserId(user.Id)
	// if p.Handle(w, r) {
	// 	return
	// }

	// TODO: Write the functionality
	// latestArticles, p := db.GetLatestArticlesByUserId(user.id)
	// if p.Handle(w, r) {
	// 	return
	// }

	socials, p := db.GetSocialPlatformsByUserId(user.Id)
	if p.Handle(w, r) {
		return
	}

	response := map[string]interface{}{
		"username":       user.Username,
		"email":          user.Email,
		"image":          user.ImageUrl,
		"bio":            user.Bio,
		"favourites":     user.Favourites,
		"created_at":     user.CreatedAt,
		"articles_count": 0,
		"words_count":    0,
		"reads_count":    0,
		"avg_rating":     0,
		"book": map[string]string{
			"title":       "",
			"description": "",
			"cover_url":   "",
		},
		"latest_articles": map[string]string{
			"title":      "",
			"tldr":       "",
			"created_at": "",
			"views":      "",
			"avg_rating": "",
		},
		"socials": []map[string]string{},
	}
	for _, social := range socials {
		s := map[string]string{
			"type":  models.SocialPlatformGetName(social.Type),
			"label": social.Label,
			"value": social.Value,
		}
		response["socials"] = append(response["socials"].([]map[string]string), s)
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		p := problems.Problem{
			Type:          problems.HandlerProblem,
			ServerMessage: "coulnd't ecnode profile data",
			ClientMessage: "An error occurred while processing your request",
			Status:        http.StatusInternalServerError,
		}
		p.Handle(w, r)
		return
	}
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	// TODO: Write the handler
}
