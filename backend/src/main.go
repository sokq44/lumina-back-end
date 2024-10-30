package main

import (
	"backend/config"
	"backend/handlers"
	"backend/middleware"
	"backend/utils/database"
	"backend/utils/emails"
	"fmt"
	"log"
	"net/http"
)

func main() {
	config.InitConfig()
	database.InitDb()
	emails.InitEmails()

	http.HandleFunc("/user/login", handlers.LoginUser)
	http.HandleFunc("/user/register", handlers.RegisterUser)
	http.HandleFunc("/user/verify-email", handlers.VerifyEmail)
	http.HandleFunc("/user/logout", middleware.Authenticate(handlers.LogoutUser))
	http.HandleFunc("/user/logged-in", middleware.Authenticate(handlers.UserLoggedIn))
	http.HandleFunc("/user/get-user", middleware.Authenticate(handlers.GetUser))
	http.HandleFunc("/user/modify-user", middleware.Authenticate(handlers.ModifyUser))

	port := config.Port

	log.Println("serving on http://localhost:"+port, "(press ctrl + c to stop the process)")
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatal("Error while trying to start the server.")
	}
}
