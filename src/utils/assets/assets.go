package assets

import (
	"backend/utils/errhandle"
	"fmt"
	"net/http"
	"os"
)

func CreateAsset(filename string, data []byte) *errhandle.Error {
	f, err := os.Create(filename)
	if err != nil {
		return &errhandle.Error{
			Type:          errhandle.HandlerError,
			ServerMessage: fmt.Sprintf("while creating the file -> %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
	}
	_, err = f.Write(data)
	if err != nil {
		return &errhandle.Error{
			Type:          errhandle.HandlerError,
			ServerMessage: fmt.Sprintf("while writing the file -> %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
	}

	if err := f.Close(); err != nil {
		return &errhandle.Error{
			Type:          errhandle.HandlerError,
			ServerMessage: fmt.Sprintf("while closing the file -> %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
	}

	return nil
}
