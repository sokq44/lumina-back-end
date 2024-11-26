package main

import (
	"backend/config"
	"backend/handlers"
	"backend/middleware"
	"backend/utils/database"
	"backend/utils/emails"
	"backend/utils/errhandle"
	"flag"
	"fmt"
	"log"
	"net/http"
)

func initApplication() {
	configPath := flag.String("config", "./.env", "Path to the [.env] configuration file.")
	logsPath := flag.String("logs", "../logs", "Path to the [logs] directory.")
	verbose := flag.Bool("verbose", false, "Should verbose to the standard output?")

	flag.Parse()

	config.InitConfig(*configPath)
	database.InitDb()
	emails.InitEmails()
	errhandle.Init(*logsPath, *verbose)
}

func initServer() {
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

	http.HandleFunc(
		"/user/password-change-init",
		middleware.Method(
			"POST",
			handlers.PasswordChangeInit,
		),
	)

	log.Println("serving on http://localhost:"+port, "(press ctrl + c to stop the process)")
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatal("Error while trying to start the server.")
	}
}

func main() {
	initApplication()
	initServer()
}
