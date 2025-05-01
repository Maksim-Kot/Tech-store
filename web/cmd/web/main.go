package main

import (
	"log"

	"github.com/Maksim-Kot/Tech-store-web/config"
	"github.com/Maksim-Kot/Tech-store-web/internal/controller/web"
	cataloggateway "github.com/Maksim-Kot/Tech-store-web/internal/gateway/catalog/http"
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

	cataloggateway := cataloggateway.New("localhost:4001")

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

	ctrl := web.New(cataloggateway, repo)

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
