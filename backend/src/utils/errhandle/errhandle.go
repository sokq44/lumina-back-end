package errhandle

import (
	"log"
	"net/http"
)

type ErrorType byte

type Error struct {
	Type          ErrorType
	ServerMessage string
	ClientMessage string
	Status        int
}

const (
	DatabaseError ErrorType = 1
	EmailsError   ErrorType = 2
	CryptError    ErrorType = 3
	JwtError      ErrorType = 4
	ModelError    ErrorType = 5
	HandlerError  ErrorType = 6
)

func (err *Error) Handle(w http.ResponseWriter) bool {
	if err == nil {
		return false
	}

	w.WriteHeader(err.Status)
	w.Write([]byte(err.ClientMessage))

	switch err.Type {
	case DatabaseError:
		log.Printf("database error -> %v ", err.ServerMessage)
	case EmailsError:
		log.Printf("emails error -> %v ", err.ServerMessage)
	case CryptError:
		log.Printf("crypt error -> %v ", err.ServerMessage)
	case JwtError:
		log.Printf("jwt error -> %v ", err.ServerMessage)
	case ModelError:
		log.Printf("model error -> %v", err.ServerMessage)
	default:
		log.Printf("unknown error type -> %v", err.ServerMessage)
	}

	return true
}
