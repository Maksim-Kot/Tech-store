package main

import (
	"log"

	"github.com/Maksim-Kot/Tech-store-orders/config"
	"github.com/Maksim-Kot/Tech-store-orders/internal/controller/orders"
	httphandler "github.com/Maksim-Kot/Tech-store-orders/internal/handler/http"
	"github.com/Maksim-Kot/Tech-store-orders/internal/repository/postgre"
	httpserver "github.com/Maksim-Kot/Tech-store-orders/internal/server/http"
)

func main() {
	cfg, err := config.New("base.yaml")
	if err != nil {
		log.Fatal(err)
	}

	repo, err := postgre.New(cfg.Database)
	if err != nil {
		log.Fatal(err)
	}
	defer repo.Close()
	log.Printf("[server] database connection pool established")

	ctrl := orders.New(repo)
	h := httphandler.New(ctrl, cfg.Api)

	srv := httpserver.New(h, cfg.Api)
	err = srv.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
