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
	ModelError    ErrorType = 5
)

func (err *Error) Handle(w http.ResponseWriter) bool {
	if err == nil {
		return false
	}

	w.WriteHeader(err.Status)

	switch err.Type {
	case DatabaseError:
		log.Printf("database error -> %v ", err.Message)
	case EmailsError:
		log.Printf("emails error -> %v ", err.Message)
	case CryptError:
		log.Printf("crypt error -> %v ", err.Message)
	case JwtError:
		log.Printf("jwt error -> %v ", err.Message)
	case ModelError:
		log.Printf("model error -> %v", err.Message)
	default:
		log.Printf("unknown error type -> %v", err.Message)
	}

	return true
}
