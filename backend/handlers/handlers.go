package handlers

import (
	"log"
	"net/http"
)

func RegisterUser(responseWriter http.ResponseWriter, request *http.Request) {
	log.Printf("REGISTER REQUEST:\n %v\n", *request)
}
