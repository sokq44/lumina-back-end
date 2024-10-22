package main

import (
	"backend/config"
	"backend/handlers"
	"backend/utils"
	"fmt"
	"log"
	"net/http"
)

func main() {
	utils.Db.Init()
	utils.Smtp.Init()

	http.HandleFunc("/user/register", handlers.RegisterUserHandler)
	http.HandleFunc("/user/verify-email", handlers.VerifyEmailHandler)
	http.HandleFunc("/user/login", handlers.LoginHandler)

	port := config.AppContext["PORT"].(string)

	log.Println("serving on http://localhost:"+port, "(press ctrl + c to stop the process)")
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatal("Error while trying to start the server.")
	}
}
