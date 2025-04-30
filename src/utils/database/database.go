package database

import (
	"backend/config"
	"backend/utils/logs"
	"backend/utils/problems"
	"database/sql"
	"fmt"
	"log"
	"net/http"
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
	cleanupInterval := config.DbCleanupInt

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

	logs.Success("Intialized the database service.")

	for range 2 {
		db.GenerateSecret()
	}

	go db.StartCleaningDb()
	go db.StartGeneratingSecrets()
}

func GetDb() *Database {
	return &db
}

func parseTime(t string) (time.Time, *problems.Problem) {
	parsed, err := time.Parse("2006-01-02 15:04:05", t)

	if err != nil {
		return time.Time{}, &problems.Problem{
			Type:          problems.DatabaseProblem,
			ServerMessage: fmt.Sprintf("while parsing datetime -> %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
	}

	return parsed, nil
}

func (db *Database) CleanDb() *problems.Problem {
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

	secrets, err := db.GetExpiredSecrets()
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

	for _, s := range secrets {
		if err := db.DeleteSecretById(s.Id); err != nil {
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

		log.Println("deleted all unverified users and hanging email verification from the database")
	}
}

func (db *Database) StartGeneratingSecrets() {
	interval := time.Duration(config.JwtSecretGenInt)
	ticker := time.NewTicker(interval)

	defer ticker.Stop()

	for range ticker.C {
		if db.GenerateSecret().Handle(nil, nil) {
			break
		}

		log.Println("generated a new jwt secret")
	}
}
