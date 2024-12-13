package config

import (
	"log"
	"os"
	"strconv"
)

var (
	Host             string
	FrontAddr        string
	AssetsPath       string
	EmailVerTime     int
	PasswdChangeTime int
	JwtAccExpTime    int
	JwtRefExpTime    int
	JwtSecretGenInt  int
	JwtSecretExpTime int
	DbUser           string
	DbPass           string
	DbNet            string
	DbHost           string
	DbPort           string
	DbName           string
	DbCleanupInt     int
	SmtpFrom         string
	SmtpUser         string
	SmtpPass         string
	SmtpHost         string
	SmtpPort         string
)

func InitConfig() {
	Host = getEnv("LUMINA_APP_HOST")
	FrontAddr = getEnv("LUMINA_APP_FRONT_ADDR")
	AssetsPath = getEnv("LUMINA_APP_ASSETS_PATH")
	EmailVerTime = getEnvInt("LUMINA_APP_EMAIL_VER_TIME")
	PasswdChangeTime = getEnvInt("LUMINA_APP_PASSWD_VER_TIME")
	JwtAccExpTime = getEnvInt("LUMINA_APP_JWT_ACCESS_EXP_TIME")
	JwtRefExpTime = getEnvInt("LUMINA_APP_JWT_REFRESH_EXP_TIME")
	JwtSecretGenInt = getEnvInt("LUMINA_APP_JWT_SECRET_GEN_INTERVAL")
	JwtSecretExpTime = getEnvInt("LUMINA_APP_JWT_SECRET_EXP_TIME")
	DbUser = getEnv("LUMINA_DB_USER")
	DbPass = getEnv("LUMINA_DB_PASSWD")
	DbNet = getEnv("LUMINA_DB_NET")
	DbHost = getEnv("LUMINA_DB_HOST")
	DbPort = getEnv("LUMINA_DB_PORT")
	DbName = getEnv("LUMINA_DB_DBNAME")
	DbCleanupInt = getEnvInt("LUMINA_DB_CLEANUP_INTERVAL")
	SmtpFrom = getEnv("LUMINA_SMTP_FROM")
	SmtpUser = getEnv("LUMINA_SMTP_USER")
	SmtpPass = getEnv("LUMINA_SMTP_PASSWD")
	SmtpHost = getEnv("LUMINA_SMTP_HOST")
	SmtpPort = getEnv("LUMINA_SMTP_PORT")
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("error while trying to get the value of %v key from the .env file", key)
	}
	return value
}

func getEnvInt(key string) int {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("error while trying to get the value of %v key from the .env file", key)
	}

	valueNumber, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("error while trying to convert string value to number: %v", err)
	}

	return valueNumber
}
