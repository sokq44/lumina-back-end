package handlers

import (
	"backend/config"
	"backend/utils/assets"
	"backend/utils/crypt"
	"backend/utils/problems"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

func AddAsset(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		p := problems.Problem{
			Type:          problems.HandlerProblem,
			ServerMessage: fmt.Sprintf("while parsing the multipart form -> %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusBadRequest,
		}
		p.Handle(w, r)
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		p := problems.Problem{
			Type:          problems.HandlerProblem,
			ServerMessage: fmt.Sprintf("while retrieving the file -> %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusBadRequest,
		}
		p.Handle(w, r)
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			p := problems.Problem{
				Type:          problems.HandlerProblem,
				ServerMessage: fmt.Sprintf("while closing the file -> %v", err),
				ClientMessage: "An error occurred while processing your request.",
				Status:        http.StatusInternalServerError,
			}
			p.Handle(w, r)
		}
	}(file)

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		p := problems.Problem{
			Type:          problems.HandlerProblem,
			ServerMessage: fmt.Sprintf("while reading the file -> %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
		p.Handle(w, r)
		return
	}

	extension := filepath.Ext(r.FormValue("filename"))
	if extension == "" {
		p := problems.Problem{
			Type:          problems.HandlerProblem,
			ServerMessage: "missing file extension",
			ClientMessage: "File extension is required.",
			Status:        http.StatusBadRequest,
		}
		p.Handle(w, r)
		return
	}

	random, p := crypt.RandomString(64)
	if p.Handle(w, r) {
		return
	}

	fileName := "/pp-" + random + extension
	if assets.CreateAsset(config.AssetsPath+fileName, fileBytes).Handle(w, r) {
		return
	}

	_, err = w.Write([]byte(config.Host + "/images" + fileName))
	if err != nil {
		p := problems.Problem{
			Type:          problems.HandlerProblem,
			ServerMessage: fmt.Sprintf("while writing the response -> %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
		p.Handle(w, r)
		return
	}
}
