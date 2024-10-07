package main

import (
	"fmt"
	"log"
	"net/http"
)

var PORT uint16 = 3000

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func main() {
	http.HandleFunc("/", HelloWorld)

	log.Printf("Serving on http://localhost:%d", PORT)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil); err != nil {
		log.Fatal("Error while trying to start the server.")
	}

}
