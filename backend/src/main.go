package main

import (
	"backend/config"
	"backend/handlers"
	"backend/utils/database"
	"backend/utils/emails"
	"backend/utils/jwt"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// serverMux := http.NewServeMux()

	config.InitConfig()
	database.InitDb()
	emails.InitEmails()

	http.HandleFunc("/user/register", handlers.RegisterUser)
	http.HandleFunc("/user/verify-email", handlers.VerifyEmail)
	http.HandleFunc("/user/login", handlers.LoginUser)
	http.HandleFunc("/user/logged-in", jwt.Middleware(handlers.UserIsLoggedIn))

	port := config.Application.PORT

	log.Println("serving on http://localhost:"+port, "(press ctrl + c to stop the process)")
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatal("Error while trying to start the server.")
	}
}
