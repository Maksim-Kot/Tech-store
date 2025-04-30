package http

import "net/http"

func (s *Server) routes() http.Handler {
	router := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	router.HandleFunc("GET /", s.handler.Home)
	router.HandleFunc("GET /catalog", s.handler.Catalog)
	router.HandleFunc("GET /category/{id}", s.handler.ProductsByCategory)
	router.HandleFunc("GET /product/{id}", s.handler.Product)

	return s.recoverPanic(logRequest(secureHeaders(router)))
}
