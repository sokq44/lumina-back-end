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

	AppContext = map[string]string{
		"PORT":       os.Getenv("APP_PORT"),
		"FRONT_ADDR": os.Getenv("APP_FRONT_ADDR"),
		"JWT_SECRET": os.Getenv("APP_JWT_SECRET"),

		"DB_USER":             os.Getenv("DB_USER"),
		"DB_PASSWD":           os.Getenv("DB_PASSWD"),
		"DB_NET":              os.Getenv("DB_NET"),
		"DB_HOST":             os.Getenv("DB_HOST"),
		"DB_PORT":             os.Getenv("DB_PORT"),
		"DB_DBNAME":           os.Getenv("DB_DBNAME"),
		"DB_CLEANUP_INTERVAL": os.Getenv("DB_CLEANUP_INTERVAL"),

		"SMTP_FROM":   os.Getenv("SMTP_FROM"),
		"SMTP_USER":   os.Getenv("SMTP_USER"),
		"SMTP_PASSWD": os.Getenv("SMTP_PASSWD"),
		"SMTP_HOST":   os.Getenv("SMTP_HOST"),
		"SMTP_PORT":   os.Getenv("SMTP_PORT"),
	}
}
