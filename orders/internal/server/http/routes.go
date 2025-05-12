package http

import "net/http"

func (s *Server) routes() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("GET /healthcheck", s.handler.HealthcheckHandler)
	router.HandleFunc("POST /order", s.handler.CreateOrderHandler)
	router.HandleFunc("GET /order/{id}", s.handler.OrderByIDHandler)
	router.HandleFunc("GET /orders/user/{id}", s.handler.OrdersByUserIDHandler)

	v1 := http.NewServeMux()
	v1.Handle("/v1/", http.StripPrefix("/v1", router))

	return v1
}
