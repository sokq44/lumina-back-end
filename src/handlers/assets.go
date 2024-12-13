package handlers

import (
	"backend/config"
	"backend/utils/assets"
	"backend/utils/crypt"
	"backend/utils/errhandle"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

func AddImage(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		e := errhandle.Error{
			Type:          errhandle.HandlerError,
			ServerMessage: fmt.Sprintf("while parsing the multipart form -> %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusBadRequest,
		}
		e.Handle(w, r)
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		e := errhandle.Error{
			Type:          errhandle.HandlerError,
			ServerMessage: fmt.Sprintf("while retrieving the file -> %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusBadRequest,
		}
		e.Handle(w, r)
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			e := errhandle.Error{
				Type:          errhandle.HandlerError,
				ServerMessage: fmt.Sprintf("while closing the file -> %v", err),
				ClientMessage: "An error occurred while processing your request.",
				Status:        http.StatusInternalServerError,
			}
			e.Handle(w, r)
		}
	}(file)

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		e := errhandle.Error{
			Type:          errhandle.HandlerError,
			ServerMessage: fmt.Sprintf("while reading the file -> %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
		e.Handle(w, r)
		return
	}

	extension := filepath.Ext(r.FormValue("filename"))
	if extension == "" {
		e := errhandle.Error{
			Type:          errhandle.HandlerError,
			ServerMessage: "missing file extension",
			ClientMessage: "File extension is required.",
			Status:        http.StatusBadRequest,
		}
		e.Handle(w, r)
		return
	}

	random, e := crypt.RandomString(64)
	if e.Handle(w, r) {
		return
	}

	fileName := "/pp-" + random + extension
	e = assets.CreateAsset(config.AssetsPath+fileName, fileBytes)
	if e.Handle(w, r) {
		return
	}

	_, err = w.Write([]byte(config.Host + "/images" + fileName))
	if err != nil {
		e := errhandle.Error{
			Type:          errhandle.HandlerError,
			ServerMessage: fmt.Sprintf("while writing the response -> %v", err),
			ClientMessage: "An error occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
		e.Handle(w, r)
		return
	}
}
