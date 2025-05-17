package http

import (
	"net/http"

	"github.com/Maksim-Kot/Tech-store-web/ui"

	"github.com/justinas/alice"
)

func (s *Server) routes() http.Handler {
	router := http.NewServeMux()

	fileServer := http.FileServer(http.FS(ui.Files))
	router.Handle("GET /static/", fileServer)

	dynamic := alice.New(s.session, s.authenticate)

	router.Handle("GET /", dynamic.ThenFunc(s.handler.Home))
	router.Handle("GET /catalog", dynamic.ThenFunc(s.handler.Catalog))
	router.Handle("GET /category/{id}", dynamic.ThenFunc(s.handler.ProductsByCategory))
	router.Handle("GET /product/{id}", dynamic.ThenFunc(s.handler.Product))

	router.Handle("GET /user/signup", dynamic.ThenFunc(s.handler.UserSignup))
	router.Handle("POST /user/signup", dynamic.ThenFunc(s.handler.UserSignupPost))
	router.Handle("GET /user/login", dynamic.ThenFunc(s.handler.UserLogin))
	router.Handle("POST /user/login", dynamic.ThenFunc(s.handler.UserLoginPost))

	router.Handle("GET /cart", dynamic.ThenFunc(s.handler.ShowCart))
	router.Handle("POST /cart/add", dynamic.ThenFunc(s.handler.AddToCart))
	router.Handle("GET /cart/remove/{id}", dynamic.ThenFunc(s.handler.RemoveFromCart))

	protected := dynamic.Append(s.requireAuthentication)

	router.Handle("GET /account/view", protected.ThenFunc(s.handler.AccountView))
	router.Handle("GET /account/orders", protected.ThenFunc(s.handler.OrdersByUser))
	router.Handle("GET /account/order/{id}", protected.ThenFunc(s.handler.Order))

	router.Handle("GET /orders/create", dynamic.ThenFunc(s.handler.CreateOrder))
	router.Handle("POST /orders/create", dynamic.ThenFunc(s.handler.CreateOrderPost))

	router.Handle("POST /user/logout", protected.ThenFunc(s.handler.UserLogoutPost))

	standard := alice.New(s.recoverPanic, logRequest, secureHeaders)

	return standard.Then(router)
}
