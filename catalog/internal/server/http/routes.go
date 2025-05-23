package http

import (
	"net/http"

	"github.com/justinas/alice"
)

func (s *Server) routes() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("GET /healthcheck", s.handler.HealthcheckHandler)
	router.HandleFunc("GET /catalog", s.handler.CategoriesHandler)
	router.HandleFunc("GET /category/{id}", s.handler.ProductsByCategoryIDHandler)
	router.HandleFunc("GET /product/{id}", s.handler.ProductByIDHandler)

	router.HandleFunc("POST /product/{id}/decrease/{amount}", s.handler.DecreaseProductQuantityHandler)
	router.HandleFunc("POST /product/{id}/increase/{amount}", s.handler.IncreaseProductQuantityHandler)

	router.HandleFunc("POST /category", s.handler.PutCategoryHandler)
	router.HandleFunc("POST /product", s.handler.PutProductHandler)

	v1 := http.NewServeMux()
	v1.Handle("/v1/", http.StripPrefix("/v1", router))

	standard := alice.New(s.recoverPanic, logRequest)

	return standard.Then(v1)
}
