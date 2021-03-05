package route

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"streamgo/core"
	"streamgo/logger"
)

type router struct {
	log logger.Logger
	service core.Service
}

func Handle(h *http.ServeMux, log logger.Logger, s core.Service) {

	o := router{log, s}

	h.HandleFunc("/upload", log.Metrics(headers(o.postUpload), "Upload"))

	h.HandleFunc("/stream", log.Metrics(headers(o.getStream), "Stream"))
}

// headers will act as middleware to give us CORS support
func headers(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next(w, r)
	}
}

func (t *router) postUpload(w http.ResponseWriter, r *http.Request) {

	t.log.Println("File Upload Endpoint Hit")

	if err := r.ParseMultipartForm(100 << 20); err != nil {
		t.log.Println("The format file parse failed.")
		jsonResponse(w, http.StatusBadRequest, "The format file parse failed.")
		return
	}

	file, handler, err := r.FormFile("clip")
	if err != nil {
		t.log.Println("The format file invalid.")
		jsonResponse(w, http.StatusBadRequest, "The format file invalid.")
		return
	}
	defer file.Close()

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	if handler.Size > 100 << 20 {
		t.log.Println("The format file size has exceeded.")
		jsonResponse(w, http.StatusBadRequest, "The format file size has exceeded.")
		return
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		t.log.Println("The format file invalid.")
		jsonResponse(w, http.StatusBadRequest, "The format file invalid.")
		return
	}

	detectedFileType := http.DetectContentType(fileBytes)
	switch detectedFileType {
	case "image/jpeg", "image/jpg":
	case "image/gif", "image/png":
	case "application/pdf":
		break
	default:
		t.log.Println("The format file invalid mimetype.")
		jsonResponse(w, http.StatusBadRequest, "The format file invalid mimetype.")
		return
	}
	fileName := randToken(12)
	fileEndings, err := mime.ExtensionsByType(detectedFileType)
	if err != nil {
		t.log.Println("The format file invalid ext.")
		jsonResponse(w, http.StatusInternalServerError, "The format file invalid ext.")
		return
	}
	newPath := filepath.Join("temp/", fileName+fileEndings[0])
	fmt.Printf("FileType: %s, File: %s\n", detectedFileType, newPath)

	newFile, err := os.Create(newPath)
	if err != nil {
		t.log.Println("The destination file create failed.")
		jsonResponse(w, http.StatusInternalServerError, "The destination file create failed.")
		return
	}
	defer newFile.Close() // idempotent, okay to call twice

	if _, err = newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
		t.log.Println("The destination file write failed.")
		jsonResponse(w, http.StatusInternalServerError, "The destination file write failed.")
		return
	}

	if err = t.service.Upload(fileName); err != nil {
		t.log.Println("The destination file write failed.")
		jsonResponse(w, http.StatusInternalServerError, "The destination file write failed.")
		return
	}

	jsonResponse(w, http.StatusOK, "")
}

func (t *router) getStream(w http.ResponseWriter, r *http.Request) {
	t.service.Stream()
}

func jsonResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprint(w, message)
}


func randToken(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
