package database

import (
	"backend/config"
	"backend/models"
	"backend/utils/errhandle"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type Database struct {
	Connection      *sql.DB
	CleanupInterval time.Duration
}

var db Database

func InitDb() {
	user := config.DbUser
	passwd := config.DbPass
	net := config.DbNet
	host := config.DbHost
	port := config.DbPort
	dbname := config.DbName
	cleanupInterval := config.DbCleanumInt

	db.CleanupInterval = time.Duration(cleanupInterval)

	dbConfig := mysql.Config{
		User:   user,
		Passwd: passwd,
		Net:    net,
		Addr:   fmt.Sprintf("%s:%s", host, port),
		DBName: dbname,
	}

	var err error
	db.Connection, err = sql.Open("mysql", dbConfig.FormatDSN())

	if err != nil {
		log.Fatalf("failed to open the connection with database: %v", err.Error())
	}

	err = db.Connection.Ping()
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err.Error())
	}

	log.Printf("intialized the database service (%v:%v)", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))

	go db.StartCleaningDb()
}

func GetDb() *Database {
	return &db
}

func parseTime(t string) (time.Time, *errhandle.Error) {
	parsed, err := time.Parse("2006-01-02 15:04:05", t)

	if err != nil {
		return time.Time{}, &errhandle.Error{
			Type:    errhandle.DatabaseError,
			Message: fmt.Sprintf("while parsing datetime -> %v", err),
			Status:  http.StatusInternalServerError,
		}
	}

	return parsed, nil
}

func (db *Database) CreateUser(u models.User) (string, *errhandle.Error) {
	id := uuid.New().String()

	_, err := db.Connection.Exec(
		"INSERT INTO users (id, username, email, password) values (?, ?, ?, ?);",
		id, u.Username, u.Email, u.Password,
	)

	if err != nil {
		return "", &errhandle.Error{
			Type:    errhandle.DatabaseError,
			Message: fmt.Sprintf("while creating a new user -> %v", err),
			Status:  http.StatusInternalServerError,
		}
	}

	return id, nil
}

func (db *Database) UpdateUser(u models.User) *errhandle.Error {
	_, err := db.Connection.Exec(
		"UPDATE users SET username=?, email=?, password=?, verified=? WHERE id=?",
		u.Username, u.Email, u.Password, u.Verified, u.Id,
	)

	if err != nil {
		return &errhandle.Error{
			Type:    errhandle.DatabaseError,
			Message: fmt.Sprintf("while updating a user -> %v", err),
			Status:  http.StatusInternalServerError,
		}
	}

	return nil
}

func (db *Database) DeleteUserById(id string) *errhandle.Error {
	_, err := db.Connection.Exec("DELETE FROM users WHERE id=?;", id)

	if err != nil {
		return &errhandle.Error{
			Type:    errhandle.DatabaseError,
			Message: fmt.Sprintf("while deleting a user by id -> %v", err),
			Status:  http.StatusInternalServerError,
		}
	}

	return nil
}

func (db *Database) GetUserById(id string) (models.User, *errhandle.Error) {
	user := models.User{Id: id}

	err := db.Connection.QueryRow(
		"SELECT username, email, verified FROM users WHERE id=?;",
		id,
	).Scan(&user.Username, &user.Email, &user.Verified)

	if err != nil {
		return models.User{}, &errhandle.Error{
			Type:    errhandle.DatabaseError,
			Message: fmt.Sprintf("error while getting a user by id: %v", err),
			Status:  http.StatusInternalServerError,
		}
	}

	return user, nil
}

func (db *Database) GetUserByEmail(email string) (models.User, *errhandle.Error) {
	var id string
	var username string
	var password string
	var verified bool

	err := db.Connection.QueryRow(
		"SELECT id, username, password, verified FROM users WHERE email=?;",
		email,
	).Scan(&id, &username, &password, &verified)

	if err != nil {
		return models.User{}, &errhandle.Error{
			Type:    errhandle.DatabaseError,
			Message: fmt.Sprintf("error while getting a user by email: %v", err),
			Status:  http.StatusInternalServerError,
		}
	}

	user := models.User{
		Id:       id,
		Username: username,
		Password: password,
		Verified: verified,
	}

	return user, nil
}

func (db *Database) UserExists(u models.User) (bool, *errhandle.Error) {
	var id string

	err := db.Connection.QueryRow(
		"SELECT id FROM users WHERE username=? or email=?;",
		u.Username, u.Email,
	).Scan(&id)

	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, &errhandle.Error{
			Type:    errhandle.DatabaseError,
			Message: fmt.Sprintf("error while checking whether a user exists: %v", err),
			Status:  http.StatusInternalServerError,
		}
	}

	return true, nil
}

func (db *Database) VerifyUser(id string) *errhandle.Error {
	_, err := db.Connection.Exec(
		"UPDATE users SET verified=TRUE WHERE id=?;",
		id,
	)

	if err != nil {
		return &errhandle.Error{
			Type:    errhandle.DatabaseError,
			Message: fmt.Sprintf("error while verifying a user: %v", err),
			Status:  http.StatusInternalServerError,
		}
	}

	return nil
}

func (db *Database) CreateEmailVerification(e models.EmailVerification) *errhandle.Error {
	_, err := db.Connection.Exec(
		"INSERT INTO email_verification (token, expires, user_id) values (?, ?, ?);",
		e.Token, e.Expires, e.UserId,
	)

	if err != nil {
		return &errhandle.Error{
			Type:    errhandle.DatabaseError,
			Message: fmt.Sprintf("error while creating an email verification: %v", err),
			Status:  http.StatusInternalServerError,
		}
	}

	return nil
}

func (db *Database) GetEmailVerificationByToken(token string) (models.EmailVerification, *errhandle.Error) {
	var id string
	var tk string
	var userId string
	var expires string

	err := db.Connection.QueryRow(
		"SELECT id, token, expires, user_id FROM email_verification WHERE token=?;",
		token,
	).Scan(&id, &tk, &expires, &userId)

	if err != nil {
		return models.EmailVerification{}, &errhandle.Error{
			Type:    errhandle.DatabaseError,
			Message: fmt.Sprintf("error while getting an email verification by token: %v", err),
			Status:  http.StatusInternalServerError,
		}
	}

	expiresTime, e := parseTime(expires)
	if e != nil {
		return models.EmailVerification{}, e
	}

	emailVerification := models.EmailVerification{
		Id:      id,
		Token:   tk,
		UserId:  userId,
		Expires: expiresTime,
	}

	return emailVerification, nil
}

func (db *Database) DeleteEmailVerificationById(id string) *errhandle.Error {
	_, err := db.Connection.Exec(
		"DELETE FROM email_verification WHERE id=?;",
		id,
	)

	if err != nil {
		return &errhandle.Error{
			Type:    errhandle.DatabaseError,
			Message: fmt.Sprintf("error while deleting an email verification by id: %v", err),
			Status:  http.StatusInternalServerError,
		}
	}

	return nil
}

func (db *Database) GetExpiredEmailVerifications() ([]models.EmailVerification, *errhandle.Error) {
	rows, err := db.Connection.Query("SELECT * FROM email_verification WHERE expires <= NOW();")

	if err != nil {
		return nil, &errhandle.Error{
			Type:    errhandle.DatabaseError,
			Message: fmt.Sprintf("error while trying to retrieve expired email verifications: %v", err),
			Status:  http.StatusInternalServerError,
		}
	}

	var expired []models.EmailVerification
	for rows.Next() {
		var verification models.EmailVerification
		var rawTime string
		if err := rows.Scan(&verification.Id, &verification.Token, &rawTime, &verification.UserId); err != nil {
			return nil, &errhandle.Error{
				Type:    errhandle.DatabaseError,
				Message: fmt.Sprintf("error while scanning expired email verifications: %v", err),
				Status:  http.StatusInternalServerError,
			}
		}

		parsed, err := parseTime(rawTime)
		if err != nil {
			return nil, err
		}

		verification.Expires = parsed
		expired = append(expired, verification)
	}

	return expired, nil
}

func (db *Database) CreateRefreshToken(token models.RefreshToken) *errhandle.Error {
	_, err := db.Connection.Exec(
		"INSERT INTO refresh_tokens (id, token, expires, user_id) values(?, ?, ?, ?)",
		token.Id, token.Token, token.Expires, token.UserId,
	)

	if err != nil {
		return &errhandle.Error{
			Type:    errhandle.DatabaseError,
			Message: fmt.Sprintf("error while creating a refresh token: %v", err),
			Status:  http.StatusInternalServerError,
		}
	}

	return nil
}

func (db *Database) GetRefreshTokenByUserId(userId string) (*models.RefreshToken, *errhandle.Error) {
	var token models.RefreshToken
	var rawTime string

	err := db.Connection.QueryRow(
		"SELECT * FROM refresh_tokens where user_id=?;",
		userId,
	).Scan(&token.Id, &token.Token, &rawTime, &token.UserId)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, &errhandle.Error{
			Type:    errhandle.DatabaseError,
			Message: fmt.Sprintf("while getting a refresh token by user id: %v", err),
			Status:  http.StatusInternalServerError,
		}
	}

	t, e := parseTime(rawTime)
	if e != nil {
		return nil, e
	}

	token.Expires = t

	return &token, nil
}

func (db *Database) DeleteRefreshTokenById(id string) *errhandle.Error {
	_, err := db.Connection.Exec(
		"DELETE FROM refresh_tokens WHERE id=?;",
		id,
	)

	if err != nil {
		return &errhandle.Error{
			Type:    errhandle.DatabaseError,
			Message: fmt.Sprintf("error while deleting a refresh token by id: %v", err),
			Status:  http.StatusInternalServerError,
		}
	}

	return nil
}

func (db *Database) DeleteRefreshTokenByToken(token string) *errhandle.Error {
	_, err := db.Connection.Exec(
		"DELETE FROM refresh_tokens WHERE token=?;",
		token,
	)

	if err != nil {
		return &errhandle.Error{
			Type:    errhandle.DatabaseError,
			Message: fmt.Sprintf("error while deleting a refresh token by token: %v", err),
			Status:  http.StatusInternalServerError,
		}
	}

	return nil
}

func (db *Database) GetExpiredRefreshTokens() ([]models.RefreshToken, *errhandle.Error) {
	rows, err := db.Connection.Query("SELECT * FROM refresh_tokens WHERE expires <= NOW();")

	if err != nil {
		return nil, &errhandle.Error{
			Type:    errhandle.DatabaseError,
			Message: fmt.Sprintf("error while trying to retrieve expired refresh tokens: %v", err),
			Status:  http.StatusInternalServerError,
		}
	}

	var expired []models.RefreshToken
	for rows.Next() {
		var token models.RefreshToken
		var rawTime string
		if err := rows.Scan(&token.Id, &token.Token, &rawTime, &token.UserId); err != nil {
			return nil, &errhandle.Error{
				Type:    errhandle.DatabaseError,
				Message: fmt.Sprintf("error while scanning expired refresh tokens: %v", err),
				Status:  http.StatusInternalServerError,
			}
		}

		parsed, err := parseTime(rawTime)
		if err != nil {
			return nil, err
		}

		token.Expires = parsed
		expired = append(expired, token)
	}

	return expired, nil
}

func (db *Database) CleanDb() *errhandle.Error {
	verifications, err := db.GetExpiredEmailVerifications()
	if err != nil {
		return err
	}

	tokens, err := db.GetExpiredRefreshTokens()
	if err != nil {
		return err
	}

	for _, v := range verifications {
		if err = db.DeleteEmailVerificationById(v.Id); err != nil {
			return err
		}
		if err = db.DeleteUserById(v.UserId); err != nil {
			return err
		}
	}

	for _, t := range tokens {
		if err = db.DeleteRefreshTokenById(t.Id); err != nil {
			return err
		}
	}

	return nil
}

func (db *Database) StartCleaningDb() {
	ticker := time.NewTicker(db.CleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		if err := db.CleanDb(); err != nil {
			log.Println(err)
			break
		}
		log.Println("deleted all unverified users and hunging email verification from the database")
	}
}
