package database

import (
	"backend/config"
	"backend/utils/errhandle"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
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
			Type:          errhandle.DatabaseError,
			ServerMessage: fmt.Sprintf("while parsing datetime -> %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
	}

	return parsed, nil
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

	passwordChanges, err := db.GetExpiredPasswordChanges()
	if err != nil {
		return err
	}

	for _, v := range verifications {
		if err := db.DeleteEmailVerificationById(v.Id); err != nil {
			return err
		}
		if err := db.DeleteUserById(v.UserId); err != nil {
			return err
		}
	}

	for _, t := range tokens {
		if err := db.DeleteRefreshTokenById(t.Id); err != nil {
			return err
		}
	}

	for _, p := range passwordChanges {
		if err := db.DeletePasswordChangeById(p.Id); err != nil {
			return err
		}
	}

	return nil
}

func (db *Database) StartCleaningDb() {
	ticker := time.NewTicker(db.CleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		if db.CleanDb().Handle(nil, nil) {
			break
		}

		log.Println("deleted all unverified users and hunging email verification from the database")
	}
}
