package http

import (
	"net/http"

	"github.com/Maksim-Kot/Tech-store-web/ui"

	"github.com/justinas/alice"
)

func (s *Server) routes() http.Handler {
	router := http.NewServeMux()

	fileServer := http.FileServer(http.FS(ui.Files))
	router.Handle("GET /static/*filepath", fileServer)

	dynamic := alice.New(s.session, s.authenticate)

	router.Handle("GET /", dynamic.ThenFunc(s.handler.Home))
	router.Handle("GET /catalog", dynamic.ThenFunc(s.handler.Catalog))
	router.Handle("GET /category/{id}", dynamic.ThenFunc(s.handler.ProductsByCategory))
	router.Handle("GET /product/{id}", dynamic.ThenFunc(s.handler.Product))

	router.Handle("GET /user/signup", dynamic.ThenFunc(s.handler.UserSignup))
	router.Handle("POST /user/signup", dynamic.ThenFunc(s.handler.UserSignupPost))
	router.Handle("GET /user/login", dynamic.ThenFunc(s.handler.UserLogin))
	router.Handle("POST /user/login", dynamic.ThenFunc(s.handler.UserLoginPost))

	protected := dynamic.Append(s.requireAuthentication)

	router.Handle("GET /account/view", protected.ThenFunc(s.handler.AccountView))
	router.Handle("POST /user/logout", protected.ThenFunc(s.handler.UserLogoutPost))

	standard := alice.New(s.recoverPanic, logRequest, secureHeaders)

	return standard.Then(router)
}
