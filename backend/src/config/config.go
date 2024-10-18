package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// TODO:
// Write an original system for retireving sensitive data. It doesn't have to
// support only [.env] files, could be replaced with some [.json] reader
// or something.

var AppContext map[string]string

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err.Error())
		return
	}

	AppContext = make(map[string]string, 14)

	AppContext["PORT"] = os.Getenv("APP_PORT")
	AppContext["FRONT_ADDR"] = os.Getenv("APP_FRONT_ADDR")

	AppContext["DB_USER"] = os.Getenv("DB_USER")
	AppContext["DB_PASSWD"] = os.Getenv("DB_PASSWD")
	AppContext["DB_NET"] = os.Getenv("DB_NET")
	AppContext["DB_HOST"] = os.Getenv("DB_HOST")
	AppContext["DB_PORT"] = os.Getenv("DB_PORT")
	AppContext["DB_DBNAME"] = os.Getenv("DB_DBNAME")
	AppContext["DB_CLEANUP_INTERVAL"] = os.Getenv("DB_CLEANUP_INTERVAL")

	AppContext["SMTP_FROM"] = os.Getenv("SMTP_FROM")
	AppContext["SMTP_USER"] = os.Getenv("SMTP_USER")
	AppContext["SMTP_PASSWD"] = os.Getenv("SMTP_PASSWD")
	AppContext["SMTP_HOST"] = os.Getenv("SMTP_HOST")
	AppContext["SMTP_PORT"] = os.Getenv("SMTP_PORT")
}
