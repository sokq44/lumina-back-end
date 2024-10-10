package main

import (
	"backend/handlers"
	"fmt"
	"log"
	"net/http"
)

var PORT uint16 = 3000

func main() {
	http.HandleFunc("/register", handlers.RegisterUser)

	log.Printf("Serving on http://localhost:%d", PORT)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil); err != nil {
		log.Fatal("Error while trying to start the server.")
	}
}
