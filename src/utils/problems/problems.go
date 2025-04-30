package problems

import (
	"backend/utils/logs"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type ProblemType byte

type Problem struct {
	Type          ProblemType
	ServerMessage string
	ClientMessage string
	Status        int
}

const (
	DatabaseProblem ProblemType = 1
	EmailsProblem   ProblemType = 2
	CryptProblem    ProblemType = 3
	JwtProblem      ProblemType = 4
	ModelProblem    ProblemType = 5
	AssetProblem    ProblemType = 6
	HandlerProblem  ProblemType = 7
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

func (e *Problem) Handle(w http.ResponseWriter, r *http.Request) bool {
	if e == nil {
		return false
	}

	if w != nil {
		w.WriteHeader(e.Status)
		_, err := w.Write([]byte(e.ClientMessage))
		if err != nil {
			logs.Error("problem while writing a response")
		}
	}

	var logMessage string
	var typeString string

	switch e.Type {
	case DatabaseProblem:
		typeString = "database problem"
	case EmailsProblem:
		typeString = "emails problem"
	case CryptProblem:
		typeString = "crypt problem"
	case JwtProblem:
		typeString = "jwt problem"
	case ModelProblem:
		typeString = "model problem"
	case AssetProblem:
		typeString = "asset problem"
	case HandlerProblem:
		typeString = "handler problem"
	default:
		typeString = "unknown problem"
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
	if !errors.Is(err, os.ErrExist) && err != nil {
		log.Println("Problem while creating the log file", err)
	} else {
		_, err = file.WriteString(fmt.Sprintf("%s\n", logMessage))
		if err != nil {
			log.Println("problem while writing a log")
		}
	}
	err = file.Close()
	if err != nil {
		log.Println("problem while closing the log file")
	}

	return true
}
