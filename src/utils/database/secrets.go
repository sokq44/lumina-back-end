package database

import (
	"backend/config"
	"backend/models"
	"backend/utils/crypt"
	"backend/utils/errhandle"
	"fmt"
	"net/http"
	"time"
)

func (db *Database) GenerateSecret() *errhandle.Error {
	randomString, err := crypt.RandomString(64)
	if err != nil {
		return err
	}

	secret := crypt.Sha256(randomString)
	expires := time.Duration(config.JwtSecretExpTime)

	_, e := db.Connection.Query(
		"INSERT INTO secrets (secret, expires) values (?, ?)",
		secret, expires,
	)

	if e != nil {
		return &errhandle.Error{
			Type:          errhandle.DatabaseError,
			ServerMessage: fmt.Sprintf("while generating a new jwt secret: %v", e),
		}
	}

	return nil
}

func (db *Database) GetLatestSecrets() ([]models.Secret, *errhandle.Error) {
	secrets := make([]models.Secret, 2)
	rawTimes := make([]string, 2)

	rows, err := db.Connection.Query("SELECT id, secret, expires FROM secrets ORDER BY expires DESC LIMIT 2;")
	if err != nil {
		return nil, &errhandle.Error{
			Type:          errhandle.DatabaseError,
			ServerMessage: fmt.Sprintf("while getting the latest jwt secret: %v", err),
			ClientMessage: "There's been an error while creating session.",
			Status:        http.StatusInternalServerError,
		}
	}
	defer rows.Close()

	for rows.Next() {
		var secret models.Secret
		var rawTime string
		if err := rows.Scan(&secret.Id, &secret.Secret, &rawTime); err != nil {
			return nil, &errhandle.Error{
				Type:          errhandle.DatabaseError,
				ServerMessage: fmt.Sprintf("while scanning jwt secrets: %v", err),
				ClientMessage: "There's been an error while creating session.",
				Status:        http.StatusInternalServerError,
			}
		}

		parsedTime, e := parseTime(rawTimes[0])
		if e != nil {
			return nil, e
		}
		secret.Expires = parsedTime

		secrets = append(secrets, secret)
	}

	if err := rows.Err(); err != nil {
		return nil, &errhandle.Error{
			Type:          errhandle.DatabaseError,
			ServerMessage: fmt.Sprintf("while iterating over jwt secrets: %v", err),
			ClientMessage: "There's been an error while creating session.",
			Status:        http.StatusInternalServerError,
		}
	}

	return secrets, nil
}
