package http

import (
	"fmt"
	"log"
	"net/http"
)

func (s *Server) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")

				s.handler.ServerErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[server] %s %s from %s (%s)", r.Method, r.URL.RequestURI(), r.RemoteAddr, r.Proto)

		next.ServeHTTP(w, r)
	})
}
