package main

import (
	"backend/config"
	"backend/handlers"
	"backend/utils/database"
	"backend/utils/emails"
	"backend/utils/logs"
	"backend/utils/problems"
	"flag"
	"fmt"
	"log"
	"net/http"
)

// TODO: handling multiple sessions from many devices

var logo = `
██╗      ██╗   ██╗ ███╗   ███╗ ██╗ ███╗   ██╗  █████╗
██║      ██║   ██║ ████╗ ████║ ██║ ████╗  ██║ ██╔══██╗
██║      ██║   ██║ ██╔████╔██║ ██║ ██╔██╗ ██║ ███████║
██║      ██║   ██║ ██║╚██╔╝██║ ██║ ██║╚██╗██║ ██╔══██║
███████╗ ╚██████╔╝ ██║ ╚═╝ ██║ ██║ ██║ ╚████║ ██║  ██║
╚══════╝  ╚═════╝  ╚═╝     ╚═╝ ╚═╝ ╚═╝  ╚═══╝ ╚═╝  ╚═╝
`

func initApplication() string {
	appPort := flag.String("p", "3000", "Port on which the application runs.")
	logsPath := flag.String("l", "./../logs", "Path to the [logs] directory.")
	verbose := flag.Bool("v", false, "Should verbose to the standard output?")
	devMode := flag.Bool("d", false, "Specifies whether the application should run in the development mode.")

	flag.Parse()

	logs.Info(logo, false)

	config.InitConfig()
	database.InitDb()
	emails.InitEmails()
	problems.Init(*logsPath, *verbose)
	handlers.InitHandlers(*devMode, *appPort)

	return *appPort
}

func main() {
	port := initApplication()

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatal("Error while trying to start the server.", err)
	}

}
