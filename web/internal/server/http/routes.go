package http

import (
	"net/http"

	"github.com/justinas/alice"
)

func (s *Server) routes() http.Handler {
	router := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	dynamic := alice.New(s.session)

	router.Handle("GET /", dynamic.ThenFunc(s.handler.Home))
	router.Handle("GET /catalog", dynamic.ThenFunc(s.handler.Catalog))
	router.Handle("GET /category/{id}", dynamic.ThenFunc(s.handler.ProductsByCategory))
	router.Handle("GET /product/{id}", dynamic.ThenFunc(s.handler.Product))

	standard := alice.New(s.recoverPanic, logRequest, secureHeaders)

	return standard.Then(router)
}
