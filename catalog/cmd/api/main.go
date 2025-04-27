package main

import (
	"log"

	"github.com/Maksim-Kot/Tech-store-catalog/config"
	"github.com/Maksim-Kot/Tech-store-catalog/internal/controller/catalog"
	httphandler "github.com/Maksim-Kot/Tech-store-catalog/internal/handler/http"
	"github.com/Maksim-Kot/Tech-store-catalog/internal/repository/postgre"
	httpserver "github.com/Maksim-Kot/Tech-store-catalog/internal/server/http"
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
	log.Printf("database connection pool established")

	ctrl := catalog.New(repo)
	h := httphandler.New(ctrl, cfg.Api)

	srv := httpserver.New(h, cfg.Api)
	err = srv.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
