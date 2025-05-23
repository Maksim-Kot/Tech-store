package http

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Maksim-Kot/Tech-store-web/config"
	httphandler "github.com/Maksim-Kot/Tech-store-web/internal/handler/http"
)

type Server struct {
	handler *httphandler.Handler
	cfg     config.APIConfig
}

func New(h *httphandler.Handler, cfg config.APIConfig) *Server {
	return &Server{h, cfg}
}

func (s *Server) Serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.cfg.Port),
		Handler:      s.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		<-quit

		log.Println("[server] shutting down server")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		log.Println("[server] completing background tasks")

		shutdownError <- nil
	}()

	log.Printf("[server] starting %s web-auth server on %s", s.cfg.Env, srv.Addr)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	log.Printf("[server] stoped web-auth server on %s", srv.Addr)

	return nil
}
