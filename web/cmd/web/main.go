package main

import (
	"log"

	"github.com/Maksim-Kot/Commons/discovery/consul"
	"github.com/Maksim-Kot/Tech-store-web/config"
	catalogcontroller "github.com/Maksim-Kot/Tech-store-web/internal/controller/catalog"
	orderscontroller "github.com/Maksim-Kot/Tech-store-web/internal/controller/orders"
	usercontroller "github.com/Maksim-Kot/Tech-store-web/internal/controller/user"
	controller "github.com/Maksim-Kot/Tech-store-web/internal/controller/web"
	cataloggateway "github.com/Maksim-Kot/Tech-store-web/internal/gateway/catalog/http"
	ordersgateway "github.com/Maksim-Kot/Tech-store-web/internal/gateway/orders/http"
	httphandler "github.com/Maksim-Kot/Tech-store-web/internal/handler/http"
	"github.com/Maksim-Kot/Tech-store-web/internal/repository/mysql"
	httpserver "github.com/Maksim-Kot/Tech-store-web/internal/server/http"
	"github.com/Maksim-Kot/Tech-store-web/internal/session"
)

func main() {
	cfg, err := config.New("base.yaml")
	if err != nil {
		log.Fatal(err)
	}

	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		log.Fatal(err)
	}

	cataloggateway := cataloggateway.New(registry)
	ordersgateway := ordersgateway.New(registry)

	repo, err := mysql.New(cfg.Database)
	if err != nil {
		log.Fatal(err)
	}
	defer repo.Close()
	log.Printf("[server] database connection pool established")

	sessionManager, err := session.New(repo.DB, cfg.Session)
	if err != nil {
		log.Fatal(err)
	}

	catalogController := catalogcontroller.New(cataloggateway)
	ordersController := orderscontroller.New(ordersgateway)
	userController := usercontroller.New(repo)

	ctrl := controller.New(catalogController, ordersController, userController)

	h, err := httphandler.New(ctrl, sessionManager)
	if err != nil {
		log.Fatal(err)
	}

	srv := httpserver.New(h, cfg.Api)
	err = srv.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
