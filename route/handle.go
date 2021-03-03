package route

import (
	"net/http"

	"streamgo/core"
	"streamgo/logger"
)

type router struct {
	service core.Service
}

func Handle(h *http.ServeMux, log logger.Logger, s core.Service) {

	o := router{s}

	h.HandleFunc("/upload", log.Metrics(headers(o.getUpload), "Upload"))

	h.HandleFunc("/stream", log.Metrics(headers(o.getStream), "Stream"))
}

// headers will act as middleware to give us CORS support
func headers(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next(w, r)
	}
}

func (t *router) getUpload(w http.ResponseWriter, r *http.Request) {

}

func (t *router) getStream(w http.ResponseWriter, r *http.Request) {

}
