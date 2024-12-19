package handlers

import (
	"log"
	"net/http"
)

func AddArticle(w http.ResponseWriter, r *http.Request) {
	log.Println(r)
}
