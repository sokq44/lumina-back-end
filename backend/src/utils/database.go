package utils

import (
	"backend/config"
	"backend/models"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type Database struct {
	Connection      *sql.DB
	CleanupInterval time.Duration
}

var Db Database

func (db *Database) Init() {
	user := config.AppContext["DB_USER"]
	passwd := config.AppContext["DB_PASSWD"]
	net := config.AppContext["DB_NET"]
	host := config.AppContext["DB_HOST"]
	port := config.AppContext["DB_PORT"]
	dbname := config.AppContext["DB_DBNAME"]

	var err error

	cleanupInterval, err := strconv.Atoi(config.AppContext["DB_CLEANUP_INTERVAL"])
	if err != nil {
		log.Fatalf("failed to convert cleanup interval from value from string to int: %v", err.Error())
	}
	db.CleanupInterval = time.Duration(cleanupInterval)

	dbConfig := mysql.Config{
		User:   user,
		Passwd: passwd,
		Net:    net,
		Addr:   fmt.Sprintf("%s:%s", host, port),
		DBName: dbname,
	}

	db.Connection, err = sql.Open("mysql", dbConfig.FormatDSN())

	if err != nil {
		log.Fatalf("failed to open the connection with database: %v", err.Error())
	}

	err = db.Connection.Ping()
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err.Error())
	}

	log.Printf("intialized the database service(%v:%v)", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))

	go db.StartHandlingUnverifiedUsers()
}

func (db *Database) CreateUser(u models.User) (string, error) {
	id := uuid.New().String()

	_, err := db.Connection.Exec("INSERT INTO users (id, username, email, password) values (?, ?, ?, ?)", id, u.Username, u.Email, u.Password)

	if err != nil {
		return "", fmt.Errorf("error while creating a user: %v", err.Error())
	}

	return id, nil
}

func (db *Database) DeleteUser(id string) error {
	_, err := db.Connection.Exec("DELETE FROM users WHERE id=?", id)

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

func (db *Database) GetUserByEmail(email string) (models.User, error) {
	var id string
	var username string
	var password string
	var verified bool

	err := db.Connection.QueryRow("SELECT id, username, password, verified FROM users WHERE email=?", email).Scan(&id, &username, &password, &verified)
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
	_, err := db.Connection.Exec("UPDATE users SET verified=TRUE WHERE id=?", userId)

	if err != nil {
		return fmt.Errorf("error while trying to verify a user: %v", err.Error())
	}

	return nil
}

func (db *Database) CreateEmailVerification(e models.EmailVerification) error {
	_, err := db.Connection.Exec("INSERT INTO email_verification (token, expires, user_id) values (?, ?, ?)", e.Token, e.Expires, e.UserId)

	if err != nil {
		return fmt.Errorf("error while creating an email_verification row: %v", err.Error())
	}

	return nil
}

func (db *Database) GetExpiredEmailVerifications() ([]models.EmailVerification, error) {
	rows, err := db.Connection.Query("SELECT * FROM email_verification WHERE expires <= NOW();")

	if err != nil {
		return nil, fmt.Errorf("error while trying to retrieve unverified email va: %v", err.Error())
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

	err := db.Connection.QueryRow("SELECT id, token, expires, user_id FROM email_verification WHERE token=?", token).Scan(&id, &tk, &expires, &userId)
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
	_, err := db.Connection.Exec("DELETE FROM email_verification WHERE id=?", id)

	if err != nil {
		return fmt.Errorf("error while trying to remove an email verification row: %v", err.Error())
	}

	return nil
}

func (db *Database) CleanDb() error {
	expired, err := db.GetExpiredEmailVerifications()
	if err != nil {
		return err
	}

	for _, e := range expired {
		if err = db.DeleteEmailVerification(e.Id); err != nil {
			return err
		}
		if err = db.DeleteUser(e.UserId); err != nil {
			return err
		}
	}

	return nil
}

func (db *Database) StartHandlingUnverifiedUsers() {
	ticker := time.NewTicker(Db.CleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		if err := db.CleanDb(); err != nil {
			log.Println(err.Error())
			break
		}
		log.Println("deleted all unverified users and hunging email verification from the database")
	}
}

func sqlDatetimeToTime(t string) (time.Time, error) {

	parsed, err := time.Parse("2006-01-02 15:04:05", t)
	if err != nil {
		return time.Time{}, fmt.Errorf("error while parsing datetime from the database: %v", err.Error())
	}

	return parsed, nil
}
