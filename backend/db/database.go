package database

import (
	"database/sql"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

func OpenDbConnection() *sql.DB {
	dbConfig := mysql.Config{
		User:   os.Getenv("DB_USER"),
		Passwd: os.Getenv("DB_PASSWD"),
		Net:    "tcp",
		Addr:   os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT"),
		DBName: os.Getenv("DB_DBNAME"),
	}

	db, err := sql.Open("mysql", dbConfig.FormatDSN())

	if err != nil {
		log.Fatalln("Failed to open the database!")
		return nil
	}

	return db
}

var DB = OpenDbConnection()
