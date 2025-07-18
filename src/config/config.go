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
	EmailChangeTime  int
	PasswdChangeTime int
	JwtAccExpTime    int
	JwtRefExpTime    int
	JwtSecretGenInt  int
	JwtSecretExpTime int
	SwaggerPath      string
	DbUser           string
	DbPass           string
	DbNet            string
	DbHost           string
	DbPort           string
	DbName           string
	DbCleanupInt     int
	AwsSesFrom       string
)

func InitConfig() {
	Host = getEnv("LUMINA_APP_HOST")
	FrontAddr = getEnv("LUMINA_APP_FRONT_ADDR")
	AssetsPath = getEnv("LUMINA_APP_ASSETS_PATH")
	EmailVerTime = getEnvInt("LUMINA_APP_EMAIL_VER_TIME")
	PasswdChangeTime = getEnvInt("LUMINA_APP_PASSWD_VER_TIME")
	JwtAccExpTime = getEnvInt("LUMINA_APP_JWT_ACCESS_EXP_TIME")
	EmailChangeTime = getEnvInt("LUMINA_APP_EMAIL_CHANGE_TIME")
	JwtRefExpTime = getEnvInt("LUMINA_APP_JWT_REFRESH_EXP_TIME")
	JwtSecretGenInt = getEnvInt("LUMINA_APP_JWT_SECRET_GEN_INTERVAL")
	JwtSecretExpTime = getEnvInt("LUMINA_APP_JWT_SECRET_EXP_TIME")
	SwaggerPath = getEnv("LUMINA_APP_SWAGGER_PATH")
	DbUser = getEnv("LUMINA_DB_USER")
	DbPass = getEnv("LUMINA_DB_PASSWD")
	DbNet = getEnv("LUMINA_DB_NET")
	DbHost = getEnv("LUMINA_DB_HOST")
	DbPort = getEnv("LUMINA_DB_PORT")
	DbName = getEnv("LUMINA_DB_DBNAME")
	DbCleanupInt = getEnvInt("LUMINA_DB_CLEANUP_INTERVAL")
	AwsSesFrom = getEnv("AWS_SES_FROM")
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
