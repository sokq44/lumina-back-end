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

var Port string
var FrontAddr string
var EmailVerTime int
var JwtSecret string
var JwtAccExpTime int
var JwtRefExpTime int

var DbUser string
var DbPass string
var DbNet string
var DbHost string
var DbPort string
var DbName string
var DbCleanumInt int

var SmtpFrom string
var SmtpUser string
var SmtpPass string
var SmtpHost string
var SmtpPort string

func InitConfig() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err.Error())
		return
	}

	Port = getEnv("APP_PORT")
	FrontAddr = getEnv("APP_FRONT_ADDR")
	EmailVerTime = getEnvInt("APP_EMAIL_VER_TIME")
	JwtSecret = getEnv("APP_JWT_SECRET")
	JwtAccExpTime = getEnvInt("APP_JWT_ACCESS_EXP_TIME")
	JwtRefExpTime = getEnvInt("APP_JWT_REFRESH_EXP_TIME")
	DbUser = getEnv("DB_USER")
	DbPass = getEnv("DB_PASSWD")
	DbNet = getEnv("DB_NET")
	DbHost = getEnv("DB_HOST")
	DbPort = getEnv("DB_PORT")
	DbName = getEnv("DB_DBNAME")
	DbCleanumInt = getEnvInt("DB_CLEANUP_INTERVAL")
	SmtpFrom = getEnv("SMTP_FROM")
	SmtpUser = getEnv("SMTP_USER")
	SmtpPass = getEnv("SMTP_PASSWD")
	SmtpHost = getEnv("SMTP_HOST")
	SmtpPort = getEnv("SMTP_PORT")
}

func getEnv(key string) string {
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
