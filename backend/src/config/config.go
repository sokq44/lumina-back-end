package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// TODO:
// Write an original system for retireving sensitive data. It doesn't have to
// support only [.env] files, could be replaced with some [.json] reader
// or something.

type ContextItem interface{}
type Context map[string]ContextItem

var AppContext Context

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err.Error())
		return
	}

	AppContext = Context{
		"PORT":           getEnvString("APP_PORT"),
		"FRONT_ADDR":     getEnvString("APP_FRONT_ADDR"),
		"EMAIL_VER_TIME": getEnvNumber("APP_EMAIL_VER_TIME"),
		"JWT_SECRET":     getEnvString("APP_JWT_SECRET"),
		"JWT_EXP_TIME":   getEnvNumber("APP_JWT_EXP_TIME"),

		"DB_USER":             getEnvString("DB_USER"),
		"DB_PASSWD":           getEnvString("DB_PASSWD"),
		"DB_NET":              getEnvString("DB_NET"),
		"DB_HOST":             getEnvString("DB_HOST"),
		"DB_PORT":             getEnvString("DB_PORT"),
		"DB_DBNAME":           getEnvString("DB_DBNAME"),
		"DB_CLEANUP_INTERVAL": getEnvNumber("DB_CLEANUP_INTERVAL"),

		"SMTP_FROM":   getEnvString("SMTP_FROM"),
		"SMTP_USER":   getEnvString("SMTP_USER"),
		"SMTP_PASSWD": getEnvString("SMTP_PASSWD"),
		"SMTP_HOST":   getEnvString("SMTP_HOST"),
		"SMTP_PORT":   getEnvString("SMTP_PORT"),
	}
}

func getEnvString(key string) ContextItem {
	value := os.Getenv(key)

	if value == "" {
		log.Fatalf("error while trying to get the value of %v key from the .env file", key)
		return nil
	}

	return value
}

func getEnvNumber(key string) ContextItem {
	value := os.Getenv(key)

	if value == "" {
		log.Fatalf("error while trying to get the value of %v key from the .env file", key)
		return nil
	}

	valueNumber, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("error while trying to convert string value to number: %v", err)
		return nil
	}

	return valueNumber
}
