package database

import (
	"backend/config"
	"backend/models"
	"database/sql"
	"errors"
	"fmt"
	"log"
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
	user := config.Application.DB_USER
	passwd := config.Application.DB_PASSWD
	net := config.Application.DB_NET
	host := config.Application.DB_HOST
	port := config.Application.DB_PORT
	dbname := config.Application.DB_DBNAME
	cleanupInterval := config.Application.DB_CLEANUP_INTERVAL

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

func sqlDatetimeToTime(t string) (time.Time, error) {

	parsed, err := time.Parse("2006-01-02 15:04:05", t)
	if err != nil {
		return time.Time{}, fmt.Errorf("error while parsing datetime from the database: %v", err.Error())
	}

	return parsed, nil
}

func (db *Database) CreateUser(u models.User) (string, error) {
	id := uuid.New().String()

	_, err := db.Connection.Exec("INSERT INTO users (id, username, email, password) values (?, ?, ?, ?);", id, u.Username, u.Email, u.Password)

	if err != nil {
		return "", fmt.Errorf("error while creating a user: %v", err.Error())
	}

	return id, nil
}

func (db *Database) DeleteUser(id string) error {
	_, err := db.Connection.Exec("DELETE FROM users WHERE id=?;", id)

	if err != nil {
		return fmt.Errorf("error while trying to delete a user: %v", err.Error())
	}

	return nil
}

func (db *Database) UserExists(u models.User) (bool, error) {
	var id string

	err := db.Connection.QueryRow("SELECT id FROM users WHERE username=? or email=?;", u.Username, u.Email).Scan(&id)

	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		log.Println(err.Error())
		return false, errors.New("error while trying to execute the query for checking whether a user exists")
	}

	return true, nil
}

func (db *Database) GetUserById(id string) (models.User, error) {
	var user models.User = models.User{
		Id: id,
	}

	err := db.Connection.QueryRow("SELECT username, email, verified FROM users WHERE id=?;", id).Scan(&user.Username, &user.Email, &user.Verified)
	if err != nil {
		return models.User{}, fmt.Errorf("error while trying to get user by id: %v", err)
	}

	return user, nil
}

func (db *Database) GetUserByEmail(email string) (models.User, error) {
	var id string
	var username string
	var password string
	var verified bool

	err := db.Connection.QueryRow("SELECT id, username, password, verified FROM users WHERE email=?;", email).Scan(&id, &username, &password, &verified)
	if err != nil {
		return models.User{}, fmt.Errorf("error while trying to get a user by email: %v", err)
	}

	user := models.User{
		Id:       id,
		Username: username,
		Password: password,
		Verified: verified,
	}

	return user, nil
}

func (db *Database) VerifyUser(userId string) error {
	_, err := db.Connection.Exec("UPDATE users SET verified=TRUE WHERE id=?;", userId)

	if err != nil {
		return fmt.Errorf("error while trying to verify a user: %v", err.Error())
	}

	return nil
}

func (db *Database) CreateEmailVerification(e models.EmailVerification) error {
	_, err := db.Connection.Exec("INSERT INTO email_verification (token, expires, user_id) values (?, ?, ?);", e.Token, e.Expires, e.UserId)

	if err != nil {
		return fmt.Errorf("error while creating an email_verification row: %v", err.Error())
	}

	return nil
}

func (db *Database) GetExpiredEmailVerifications() ([]models.EmailVerification, error) {
	rows, err := db.Connection.Query("SELECT * FROM email_verification WHERE expires <= NOW();")

	if err != nil {
		return nil, fmt.Errorf("error while trying to retrieve unverified emails: %v", err.Error())
	}

	var unverified []models.EmailVerification
	for rows.Next() {
		var verification models.EmailVerification
		var rawTime string
		if err := rows.Scan(&verification.Id, &verification.Token, &rawTime, &verification.UserId); err != nil {
			return nil, fmt.Errorf("error while trying to scan from one of the retrieved unverified email verifications: %v", err.Error())
		}

		parsed, err := sqlDatetimeToTime(rawTime)
		if err != nil {
			return nil, err
		}

		verification.Expires = parsed
		unverified = append(unverified, verification)
	}

	return unverified, nil
}

func (db *Database) GetEmailVerificationFromToken(token string) (models.EmailVerification, error) {
	var id string
	var tk string
	var userId string
	var expires string

	err := db.Connection.QueryRow("SELECT id, token, expires, user_id FROM email_verification WHERE token=?;", token).Scan(&id, &tk, &expires, &userId)
	if err != nil {
		return models.EmailVerification{}, fmt.Errorf("error while retrieving email verification data: %v", err.Error())
	}

	expiresTime, err := sqlDatetimeToTime(expires)
	if err != nil {
		return models.EmailVerification{}, err
	}

	emailVerification := models.EmailVerification{
		Id:      id,
		Token:   tk,
		UserId:  userId,
		Expires: expiresTime,
	}

	return emailVerification, nil
}

func (db *Database) DeleteEmailVerification(id string) error {
	_, err := db.Connection.Exec("DELETE FROM email_verification WHERE id=?;", id)

	if err != nil {
		return fmt.Errorf("error while trying to remove an email verification row: %v", err.Error())
	}

	return nil
}

func (db *Database) CreateRefreshToken(token models.RefreshToken) error {
	_, err := db.Connection.Exec("INSERT INTO refresh_tokens (id, token, expires, user_id) values(?, ?, ?, ?)", token.Id, token.Token, token.Expires, token.UserId)

	if err != nil {
		return fmt.Errorf("error while trying to store a refresh token in the db: %v", err)
	}

	return nil
}

func (db *Database) GetRefreshTokenByUserId(userId string) (*models.RefreshToken, error) {
	var token models.RefreshToken
	var rawTime string

	err := db.Connection.QueryRow("SELECT * FROM refresh_tokens where user_id=?;", userId).Scan(&token.Id, &token.Token, &rawTime, &token.UserId)
	if err != nil {
		return nil, err
	}

	t, err := sqlDatetimeToTime(rawTime)
	if err != nil {
		return nil, err
	}

	token.Expires = t

	return &token, nil
}

func (db *Database) DeleteRefreshTokenById(id string) error {
	_, err := db.Connection.Exec("DELETE FROM refresh_tokens WHERE id=? OR token=?;", id)

	if err != nil {
		return fmt.Errorf("error whilte trying to delete a refresh token from the db: %v", err)
	}

	return nil
}

func (db *Database) DeleteRefreshTokenByToken(token string) error {
	_, err := db.Connection.Exec("DELETE FROM refresh_tokens WHERE token=?;", token)

	if err != nil {
		return fmt.Errorf("error whilte trying to delete a refresh token from the db: %v", err)
	}

	return nil
}

func (db *Database) GetExpiredRefreshTokens() ([]models.RefreshToken, error) {
	rows, err := db.Connection.Query("SELECT * FROM refresh_tokens WHERE expires <= NOW();")

	if err != nil {
		return nil, fmt.Errorf("error while trying to retireve expired refresh tokens: %v", err)
	}

	var expired []models.RefreshToken
	for rows.Next() {
		var token models.RefreshToken
		var rawTime string
		if err := rows.Scan(&token.Id, &token.Token, &rawTime, &token.UserId); err != nil {
			return nil, fmt.Errorf("error while trying to scan from one of the expired refresh tokens: %v", err)
		}

		parsed, err := sqlDatetimeToTime(rawTime)
		if err != nil {
			return nil, err
		}

		token.Expires = parsed
		expired = append(expired, token)
	}

	return expired, nil
}

func (db *Database) CleanDb() error {
	verifications, err := db.GetExpiredEmailVerifications()
	if err != nil {
		return err
	}

	tokens, err := db.GetExpiredRefreshTokens()
	if err != nil {
		return err
	}

	for _, v := range verifications {
		if err = db.DeleteEmailVerification(v.Id); err != nil {
			return err
		}
		if err = db.DeleteUser(v.UserId); err != nil {
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
			log.Println(err.Error())
			break
		}
		log.Println("deleted all unverified users and hunging email verification from the database")
	}
}
