package http

import (
	"log"
	"net/http"
)

func (h *Handler) logError(r *http.Request, err error) {
	log.Printf("url: %s %v", r.URL, err)
}

func (h *Handler) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"error": message}

	err := h.writeJSON(w, status, env, nil)
	if err != nil {
		h.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	h.logError(r, err)
	message := "the server encountered a problem and could not process your request"
	h.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (h *Handler) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	h.errorResponse(w, r, http.StatusNotFound, message)
}

func (h *Handler) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	h.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (h *Handler) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to update due to an edit conflict, please try again"
	h.errorResponse(w, r, http.StatusConflict, message)
}
