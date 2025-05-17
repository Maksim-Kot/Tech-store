package http

import (
	"net/http"

	"github.com/justinas/alice"
)

func (s *Server) routes() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("GET /healthcheck", s.handler.HealthcheckHandler)
	router.HandleFunc("POST /order", s.handler.CreateOrderHandler)
	router.HandleFunc("GET /order/{id}", s.handler.OrderByIDHandler)
	router.HandleFunc("GET /orders/user/{id}", s.handler.OrdersByUserIDHandler)

	v1 := http.NewServeMux()
	v1.Handle("/v1/", http.StripPrefix("/v1", router))

	standard := alice.New(s.recoverPanic, logRequest)

	return standard.Then(v1)
}
