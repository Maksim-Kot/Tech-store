package main

import (
	"log"

	"github.com/Maksim-Kot/Tech-store-web/config"
	"github.com/Maksim-Kot/Tech-store-web/internal/controller/web"
	cataloggateway "github.com/Maksim-Kot/Tech-store-web/internal/gateway/catalog/http"
	httphandler "github.com/Maksim-Kot/Tech-store-web/internal/handler/http"
	httpserver "github.com/Maksim-Kot/Tech-store-web/internal/server/http"
)

func main() {
	cfg, err := config.New("base.yaml")
	if err != nil {
		log.Fatal(err)
	}

	cataloggateway := cataloggateway.New("localhost:4001")

	ctrl := web.New(cataloggateway)
	h, err := httphandler.New(ctrl)
	if err != nil {
		log.Fatal(err)
	}

	srv := httpserver.New(h, cfg.Api)
	err = srv.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
