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

func initApplication() string {
	appPort := flag.String("p", "3000", "Port on which the application runs.")
	logsPath := flag.String("l", "./../logs", "Path to the [logs] directory.")
	verbose := flag.Bool("v", false, "Should verbose to the standard output?")

	flag.Parse()

	config.InitConfig()
	database.InitDb()
	emails.InitEmails()
	errhandle.Init(*logsPath, *verbose)

	return *appPort
}

func initServer(port string) {
	http.HandleFunc(
		"/user/login",
		middleware.CORS(
			middleware.Method(
				"POST",
				handlers.LoginUser,
			),
		),
	)

	http.HandleFunc(
		"/user/register",
		middleware.CORS(
			middleware.Method(
				"POST",
				handlers.RegisterUser,
			),
		),
	)

	http.HandleFunc(
		"/user/verify-email",
		middleware.CORS(
			middleware.Method(
				"PATCH",
				handlers.VerifyEmail,
			),
		),
	)

	http.HandleFunc(
		"/user/logout",
		middleware.CORS(
			middleware.Authenticate(
				middleware.Method(
					"DELETE",
					handlers.LogoutUser,
				),
			),
		),
	)

	http.HandleFunc(
		"/user/logged-in",
		middleware.CORS(
			middleware.Authenticate(
				middleware.Method(
					"GET",
					func(w http.ResponseWriter, r *http.Request) {},
				),
			),
		),
	)

	http.HandleFunc(
		"/user/get-user",
		middleware.CORS(
			middleware.Authenticate(
				middleware.Method(
					"GET",
					handlers.GetUser,
				),
			),
		),
	)

	http.HandleFunc(
		"/user/modify-user",
		middleware.CORS(
			middleware.Authenticate(
				middleware.Method(
					"PATCH",
					handlers.ModifyUser,
				),
			),
		),
	)

	http.HandleFunc(
		"/user/change-password",
		middleware.CORS(
			middleware.Method(
				"PATCH",
				handlers.ChangePassword,
			),
		),
	)

	http.HandleFunc(
		"/user/password-change-init",
		middleware.CORS(
			middleware.Method(
				"POST",
				handlers.PasswordChangeInit,
			),
		),
	)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatal("Error while trying to start the server.")
	}
}

func main() {
	port := initApplication()
	initServer(port)
}
