package utils

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
	Connection *sql.DB
}

var Db Database

func init() {
	user := config.AppContext["DB_USER"]
	passwd := config.AppContext["DB_PASSWD"]
	net := config.AppContext["DB_NET"]
	host := config.AppContext["DB_HOST"]
	port := config.AppContext["DB_PORT"]
	dbname := config.AppContext["DB_DBNAME"]

	dbConfig := mysql.Config{
		User:   user,
		Passwd: passwd,
		Net:    net,
		Addr:   fmt.Sprintf("%s:%s", host, port),
		DBName: dbname,
	}

	var err error
	Db.Connection, err = sql.Open("mysql", dbConfig.FormatDSN())

	if err != nil {
		log.Fatalf("failed to open the connection with database: %v", err.Error())
	}

	err = Db.Connection.Ping()
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err.Error())
	}

	log.Printf("connected to dbms server: %v:%v", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))
}

func (db *Database) CreateUser(u models.User) (string, error) {
	id := uuid.New().String()

	_, err := db.Connection.Exec("INSERT INTO users (id, username, email, password) values (?, ?, ?, ?)", id, u.Username, u.Email, u.Password)

	if err != nil {
		return "", fmt.Errorf("error while creating a user: %v", err.Error())
	}

	return id, nil
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

func (db *Database) VerifyUser(userId string) error {
	_, err := db.Connection.Exec("UPDATE users SET verified=TRUE WHERE id=?", userId)

	if err != nil {
		return fmt.Errorf("error while trying to verify a user: %v", err.Error())
	}

	return nil
}

// FIXME:
// This method should take models.EmailValidation as an argument.
func (db *Database) CreateEmailValidation(userId, token string, expires time.Time) error {
	_, err := db.Connection.Exec("INSERT INTO email_verification (token, expires, user_id) values (?, ?, ?)", token, expires, userId)

	if err != nil {
		return fmt.Errorf("error while creating an email_verification row: %v", err.Error())
	}

	return nil
}

func (db *Database) GetEmailValidationFromToken(token string) (models.EmailValidation, error) {
	var id string
	var tk string
	var userId string
	var expires string

	err := db.Connection.QueryRow("SELECT id, token, expires, user_id FROM email_verification WHERE token=?", token).Scan(&id, &tk, &expires, &userId)
	if err != nil {
		return models.EmailValidation{}, fmt.Errorf("error while retrieving email validation data: %v", err.Error())
	}

	const layout = "2006-01-02 15:04:05"
	expiresTime, err := time.Parse(layout, expires)
	if err != nil {
		return models.EmailValidation{}, fmt.Errorf("error while parsing datetime from the database: %v", err.Error())
	}

	emailValidation := models.EmailValidation{
		Id:      id,
		Token:   tk,
		UserId:  userId,
		Expires: expiresTime,
	}

	return emailValidation, nil
}

func (db *Database) DeleteEmailValidation(id string) error {
	_, err := db.Connection.Exec("DELETE FROM email_verification WHERE id=?", id)

	if err != nil {
		return fmt.Errorf("error while trying to remove an email validation row: %v", err.Error())
	}

	return nil
}
