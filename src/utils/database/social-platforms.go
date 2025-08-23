package database

import (
	"backend/models"
	"backend/utils/problems"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
)

func (db *Database) GetSocialPlatforms() ([]models.SocialPlatform, *problems.Problem) {

	q := "SELECT type, name FORM socials;"

	rows, err := db.Connection.Query(q)
	if err != nil {
		rows.Close()
		return nil, &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("querying socials failed: %v", err),
			ClientMessage: "An unexpected error has occurred when fetching the available social platforms.",
			Status:        http.StatusInternalServerError,
		}
	}
	defer rows.Close()

	socials := make([]models.SocialPlatform, 0)
	for rows.Next() {
		var s models.SocialPlatform
		if err = rows.Scan(&s.Type, &s.Name); err != nil {
			return nil, &problems.Problem{
				Type:          problems.DatabaseProblem,
				ServerMessage: fmt.Sprintf("querying socials failed: %v", err),
				ClientMessage: "An unexpected error has occurred when fetching the available social platforms.",
				Status:        http.StatusInternalServerError,
			}
		}
		socials = append(socials, s)
	}

	rows.Close()

	return socials, nil
}

func (db *Database) GetSocialPlatformsByUserId(userId string) ([]models.UserSocialPlatform, *problems.Problem) {
	q := "SELECT social_type, value, label FROM user_socials WHERE user_id=?;"
	rows, err := db.Connection.Query(q, userId)
	if errors.Is(err, sql.ErrNoRows) {
		rows.Close()
		return nil, &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("no social platforms were found -> %v", err),
			ClientMessage: "An unexpected error error has occured when trying to fetch all the available social platforms.",
			Status:        http.StatusNotFound,
		}
	} else if err != nil {
		rows.Close()
		return nil, &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("while trying to get all the available social platforms -> %v", err),
			ClientMessage: "An unexpected error error has occured when trying to fetch all the available social platforms.",
			Status:        http.StatusInternalServerError,
		}
	}
	defer rows.Close()

	socials := make([]models.UserSocialPlatform, 0)
	for rows.Next() {
		s := models.UserSocialPlatform{UserId: userId}
		if err = rows.Scan(&s.Type, &s.Value, &s.Label); err != nil {
			return nil, &problems.Problem{
				Type:          problems.DatabaseProblem,
				ServerMessage: fmt.Sprintf("querying socials failed: %v", err),
				ClientMessage: "An unexpected error has occurred when fetching the available social platforms.",
				Status:        http.StatusInternalServerError,
			}
		}
		socials = append(socials, s)
	}

	return socials, nil
}
