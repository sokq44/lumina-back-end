package main

import (
	"backend/config"
	"backend/handlers"
	"backend/utils/database"
	"backend/utils/emails"
	"backend/utils/problems"
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
	handlers.InitHandlers()
	problems.Init(*logsPath, *verbose)

	return *appPort
}

func initServer(port string) {
	initApplication()

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatal("Error while trying to start the server.")
	}
}

func main() {
	port := initApplication()
	initServer(port)
}
