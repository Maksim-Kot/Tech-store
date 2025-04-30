package http

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
)

func (h *Handler) ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	log.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (h *Handler) ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (h *Handler) NotFound(w http.ResponseWriter) {
	h.ClientError(w, http.StatusNotFound)
}
