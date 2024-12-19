package assets

import (
	"backend/utils/problems"
	"fmt"
	"net/http"
	"os"
)

func CreateAsset(filename string, data []byte) *problems.Problem {
	f, err := os.Create(filename)
	if err != nil {
		return &problems.Problem{
			Type:          problems.AssetProblem,
			ServerMessage: fmt.Sprintf("while creating the file -> %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
	}
	_, err = f.Write(data)
	if err != nil {
		return &problems.Problem{
			Type:          problems.AssetProblem,
			ServerMessage: fmt.Sprintf("while writing the file -> %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
	}

	if err := f.Close(); err != nil {
		return &problems.Problem{
			Type:          problems.AssetProblem,
			ServerMessage: fmt.Sprintf("while closing the file -> %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
	}

	return nil
}
