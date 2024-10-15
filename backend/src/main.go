package main

import (
	database "backend/db"
	user "backend/handlers"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var PORT uint16 = 3000

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	err = database.OpenDbConnection(os.Getenv("DB_USER"), os.Getenv("DB_PASSWD"), "tcp", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_DBNAME"))
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	http.HandleFunc("/register", user.RegisterUserHandler)

	log.Printf("Serving on http://localhost:%d", PORT)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil); err != nil {
		log.Fatal("Error while trying to start the server.")
	}
}
