package web

import (
	catalogcontroller "github.com/Maksim-Kot/Tech-store-web/internal/controller/catalog"
	orderscontroller "github.com/Maksim-Kot/Tech-store-web/internal/controller/orders"
	usercontroller "github.com/Maksim-Kot/Tech-store-web/internal/controller/user"
)

type Controller struct {
	Catalog *catalogcontroller.CatalogController
	Orders  *orderscontroller.OrdersController
	User    *usercontroller.UserController
}

func New(catalog *catalogcontroller.CatalogController, orders *orderscontroller.OrdersController, user *usercontroller.UserController) *Controller {
	return &Controller{
		Catalog: catalog,
		Orders:  orders,
		User:    user,
	}
}
