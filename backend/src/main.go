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
	utils.JWT.Init()
	utils.Smtp.Init()
	utils.Crypto.Init()

	http.HandleFunc("/user/register", handlers.RegisterUserHandler)
	http.HandleFunc("/user/verify-email", handlers.VerifyEmailHandler)

	log.Println("serving on http://localhost:"+config.AppContext["PORT"], "(press ctrl + c to stop the process)")
	if err := http.ListenAndServe(fmt.Sprintf(":%s", config.AppContext["PORT"]), nil); err != nil {
		log.Fatal("Error while trying to start the server.")
	}
}
