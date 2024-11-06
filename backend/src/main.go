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

	port := config.Port

	http.HandleFunc(
		"/user/login",
		middleware.Method(
			"POST",
			handlers.LoginUser,
		),
	)

	http.HandleFunc(
		"/user/register",
		middleware.Method(
			"POST",
			handlers.RegisterUser,
		),
	)

	http.HandleFunc(
		"/user/verify-email",
		middleware.Method(
			"PATCH",
			handlers.VerifyEmail,
		),
	)

	http.HandleFunc(
		"/user/logout",
		middleware.Authenticate(
			middleware.Method(
				"DELETE",
				handlers.LogoutUser,
			),
		),
	)

	http.HandleFunc(
		"/user/logged-in",
		middleware.Authenticate(
			middleware.Method(
				"GET",
				func(w http.ResponseWriter, r *http.Request) {},
			),
		),
	)

	http.HandleFunc(
		"/user/get-user",
		middleware.Authenticate(
			middleware.Method(
				"GET",
				handlers.GetUser,
			),
		),
	)

	http.HandleFunc(
		"/user/modify-user",
		middleware.Authenticate(
			middleware.Method(
				"PATCH",
				handlers.ModifyUser,
			),
		),
	)

	http.HandleFunc(
		"/user/change-password",
		middleware.Method(
			"PATCH",
			handlers.ChangePassword,
		),
	)

	log.Println("serving on http://localhost:"+port, "(press ctrl + c to stop the process)")
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatal("Error while trying to start the server.")
	}
}
