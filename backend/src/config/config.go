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

type Context struct {
	PORT                 string
	FRONT_ADDR           string
	EMAIL_VER_TIME       int
	JWT_SECRET           string
	JWT_ACCESS_EXP_TIME  int
	JWT_REFRESH_EXP_TIME int

	DB_USER             string
	DB_PASSWD           string
	DB_NET              string
	DB_HOST             string
	DB_PORT             string
	DB_DBNAME           string
	DB_CLEANUP_INTERVAL int

	SMTP_FROM   string
	SMTP_USER   string
	SMTP_PASSWD string
	SMTP_HOST   string
	SMTP_PORT   string
}

var Application Context

func InitConfig() {
	if err := godotenv.Load("config/.env"); err != nil {
		log.Fatal(err.Error())
		return
	}

	Application = Context{
		PORT:                 getEnvString("APP_PORT"),
		FRONT_ADDR:           getEnvString("APP_FRONT_ADDR"),
		EMAIL_VER_TIME:       getEnvInt("APP_EMAIL_VER_TIME"),
		JWT_SECRET:           getEnvString("APP_JWT_SECRET"),
		JWT_ACCESS_EXP_TIME:  getEnvInt("APP_JWT_ACCESS_EXP_TIME"),
		JWT_REFRESH_EXP_TIME: getEnvInt("APP_JWT_REFRESH_EXP_TIME"),
		DB_USER:              getEnvString("DB_USER"),
		DB_PASSWD:            getEnvString("DB_PASSWD"),
		DB_NET:               getEnvString("DB_NET"),
		DB_HOST:              getEnvString("DB_HOST"),
		DB_PORT:              getEnvString("DB_PORT"),
		DB_DBNAME:            getEnvString("DB_DBNAME"),
		DB_CLEANUP_INTERVAL:  getEnvInt("DB_CLEANUP_INTERVAL"),
		SMTP_FROM:            getEnvString("SMTP_FROM"),
		SMTP_USER:            getEnvString("SMTP_USER"),
		SMTP_PASSWD:          getEnvString("SMTP_PASSWD"),
		SMTP_HOST:            getEnvString("SMTP_HOST"),
		SMTP_PORT:            getEnvString("SMTP_PORT"),
	}
}

func getEnvString(key string) string {
	value := os.Getenv(key)

	if value == "" {
		log.Fatalf("error while trying to get the value of %v key from the .env file", key)
		return ""
	}

	return value
}

func getEnvInt(key string) int {
	value := os.Getenv(key)

	if value == "" {
		log.Fatalf("error while trying to get the value of %v key from the .env file", key)
		return -1
	}

	valueNumber, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("error while trying to convert string value to number: %v", err)
		return -1
	}

	return valueNumber
}
