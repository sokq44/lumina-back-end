package errhandle

import (
	"log"
	"net/http"
)

type ErrorType byte

type Error struct {
	Type    ErrorType
	Message string
	Status  int
}

const (
	DatabaseError ErrorType = 1
	EmailsError   ErrorType = 2
	CryptError    ErrorType = 3
	JwtError      ErrorType = 4
)

func (err *Error) Handle(w http.ResponseWriter) bool {
	if err == nil {
		return false
	}

	if err.Type == DatabaseError {
		w.WriteHeader(err.Status)
		log.Printf("database error -> %v ", err.Message)
	} else if err.Type == EmailsError {
		w.WriteHeader(err.Status)
		log.Printf("emails error -> %v ", err.Message)
	} else if err.Type == CryptError {
		w.WriteHeader(err.Status)
		log.Printf("crypt error -> %v ", err.Message)
	} else if err.Type == JwtError {
		w.WriteHeader(err.Status)
		log.Printf("jwt error -> %v ", err.Message)
	} else {
		log.Printf("can't handle an unknown type of error: %v", err.Type)
		return false
	}

	return true
}
