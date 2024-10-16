package main

import (
	"backend/handlers"
	"backend/utils"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err.Error())
		return
	}

	if msg, err := utils.Db.OpenDbConnection(
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWD"),
		os.Getenv("DB_NET"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DBNAME"),
	); err != nil {
		log.Fatal(err.Error())
		return
	} else {
		log.Println(msg)
	}

	if msg, err := utils.Smtp.OpenSmtpConnection(
		os.Getenv("SMTP_FROM"),
		os.Getenv("SMTP_USER"),
		os.Getenv("SMTP_PASSWD"),
		os.Getenv("SMTP_HOST"),
		os.Getenv("SMTP_PORT"),
	); err != nil {
		log.Fatal(err.Error())
		return
	} else {
		log.Println(msg)
	}

	port := os.Getenv("APP_PORT")

	http.HandleFunc("/user/register", handlers.RegisterUserHandler)

	log.Println("serving on http://localhost:" + port)
	if err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil); err != nil {
		log.Fatal("Error while trying to start the server.")
	}
}
