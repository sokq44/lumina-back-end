package main

import (
	"backend/config"
	"backend/handlers"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/user/register", handlers.RegisterUserHandler)
	http.HandleFunc("/user/verify-email", handlers.VerifyEmailHandler)

	log.Println("serving on http://localhost:"+config.AppContext["PORT"], "(press ctrl + c to stop the process)")
	if err := http.ListenAndServe(fmt.Sprintf(":%s", config.AppContext["PORT"]), nil); err != nil {
		log.Fatal("Error while trying to start the server.")
	}
}
