package errhandle

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
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

var location string
var shouldVerbose bool

func Init(l string, v bool) {
	location = l
	shouldVerbose = v

	err := os.MkdirAll(location, os.ModePerm)
	if err != nil {
		log.Fatalf("failed to create the logs directory")
	}
}

func (e *Error) Handle(w http.ResponseWriter, r *http.Request) bool {
	if e == nil {
		return false
	}

	if w != nil {
		w.WriteHeader(e.Status)
		w.Write([]byte(e.ClientMessage))
	}

	var logMessage string
	var typeString string

	switch e.Type {
	case DatabaseError:
		typeString = "database error"
	case EmailsError:
		typeString = "emails error"
	case CryptError:
		typeString = "crypt error"
	case JwtError:
		typeString = "jwt error"
	case ModelError:
		typeString = "model error"
	default:
		typeString = "unknown error"
	}

	logMessage = fmt.Sprintf(
		"\n{\n\t%s -> %s,\n\tresponded with HTTP status -> %d\n\tresponded with message -> \"%s\"\n}\n",
		typeString, e.ServerMessage, e.Status, e.ClientMessage,
	)

	var host string
	if r != nil {
		host = r.Host
	}
	now := time.Now()
	logMessage = fmt.Sprintf("[%v] [%s] -> %s", now, host, logMessage)

	if shouldVerbose {
		log.Println(logMessage)
	}

	fullPath := filepath.Join(location, now.Format("02-01-2006")+".log")
	file, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != os.ErrExist && err != nil {
		log.Println("error while creating the log file", err)
	} else {
		_, err = file.WriteString(fmt.Sprintf("%s\n", logMessage))
		if err != nil {
			log.Println("error while writing a log")
		}
	}
	file.Close()

	return true
}
