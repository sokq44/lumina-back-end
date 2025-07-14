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
	verifications, p := db.GetExpiredEmailVerifications()
	if p != nil {
		return p
	}

	tokens, p := db.GetExpiredRefreshTokens()
	if p != nil {
		return p
	}

	emailChanges, p := db.GetExpiredEmailChanges()
	if p != nil {
		return p
	}

	passwordChanges, p := db.GetExpiredPasswordChanges()
	if p != nil {
		return p
	}

	secrets, p := db.GetExpiredSecrets()
	if p != nil {
		return p
	}

	for _, v := range verifications {
		if p := db.DeleteEmailVerificationById(v.Id); p != nil {
			return p
		}
		if p := db.DeleteUserById(v.UserId); p != nil {
			return p
		}
	}

	for _, t := range tokens {
		if p := db.DeleteRefreshTokenById(t.Id); p != nil {
			return p
		}
	}

	for _, ec := range emailChanges {
		if p := db.DeleteEmailChangeById(ec.Id); p != nil {
			return p
		}
	}

	for _, pc := range passwordChanges {
		if p := db.DeletePasswordChangeById(pc.Id); p != nil {
			return p
		}
	}

	for _, s := range secrets {
		if p := db.DeleteSecretById(s.Id); p != nil {
			return p
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

		logs.Info("deleted all unverified users and hanging email verification from the database")
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

		logs.Info("generated a new jwt secret")
	}
}
