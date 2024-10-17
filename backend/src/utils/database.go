package utils

import (
	"backend/models"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type Database struct {
	Connection *sql.DB
	User       string
	Passwd     string
	Net        string
	Host       string
	Port       string
	DBName     string
}

var Db Database

func (db *Database) OpenDbConnection(user string, passwd string, net string, host string, port string, dbname string) (string, error) {
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
		return "", fmt.Errorf("failed to open the connection with database: %v", err.Error())
	}

	err = db.Connection.Ping()
	if err != nil {
		return "", fmt.Errorf("failed to connect to the database: %v", err.Error())
	}

	db.User = user
	db.Passwd = passwd
	db.Net = net
	db.Host = host
	db.Port = port
	db.DBName = dbname

	return fmt.Sprintf("connected to dbms server: %v:%v", host, port), nil
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

func (db *Database) CreateEmailValidation(userId, token string, expires time.Time) error {
	_, err := db.Connection.Exec("INSERT INTO email_validation (token, expires, user_id) values (?, ?, ?)", token, expires, userId)

	if err != nil {
		return fmt.Errorf("error while creating an email_validation row: %v", err.Error())
	}

	return nil
}
