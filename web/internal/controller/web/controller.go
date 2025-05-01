package web

import (
	catalogcontroller "github.com/Maksim-Kot/Tech-store-web/internal/controller/catalog"
	usercontroller "github.com/Maksim-Kot/Tech-store-web/internal/controller/user"
)

type Controller struct {
	Catalog *catalogcontroller.CatalogController
	User    *usercontroller.UserController
}

func New(catalog *catalogcontroller.CatalogController, user *usercontroller.UserController) *Controller {
	return &Controller{
		Catalog: catalog,
		User:    user,
	}
}
