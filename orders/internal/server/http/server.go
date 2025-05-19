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

	"github.com/Maksim-Kot/Commons/discovery"
	"github.com/Maksim-Kot/Tech-store-orders/config"
	httphandler "github.com/Maksim-Kot/Tech-store-orders/internal/handler/http"
)

type Server struct {
	handler    *httphandler.Handler
	cfg        config.APIConfig
	registry   discovery.Registry
	instanceID string
}

func New(h *httphandler.Handler, cfg config.APIConfig, registry discovery.Registry) *Server {
	return &Server{
		handler:  h,
		cfg:      cfg,
		registry: registry,
	}
}

func (s *Server) Serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.cfg.Port),
		Handler:      s.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	ctx, mainCancel := context.WithCancel(context.Background())
	defer mainCancel()

	if err := s.registerService(ctx); err != nil {
		return fmt.Errorf("failed to register service: %w", err)
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

		mainCancel()
		s.deregisterService(ctx)

		shutdownError <- nil
	}()

	log.Printf("[server] starting %s orders server on %s", s.cfg.Env, srv.Addr)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	log.Printf("[server] stoped orders server on %s", srv.Addr)

	return nil
}

func (s *Server) registerService(ctx context.Context) error {
	s.instanceID = discovery.GenerateInstanceID(s.cfg.Name)
	addr := fmt.Sprintf("localhost:%d", s.cfg.Port)

	if err := s.registry.Register(ctx, s.instanceID, s.cfg.Name, addr); err != nil {
		return err
	}

	go s.reportHealth(ctx)

	return nil
}

func (s *Server) deregisterService(ctx context.Context) {
	log.Println("[registry] deregistering service from Consul...")
	if err := s.registry.Deregister(ctx, s.instanceID, s.cfg.Name); err != nil {
		log.Println("[registry] failed to deregister from Consul:", err)
	} else {
		log.Println("[registry] successfully deregistered from Consul")
	}
}

func (s *Server) reportHealth(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("[registry] stopping health reporting")
			return
		default:
			if err := s.registry.ReportHealthyState(s.instanceID, s.cfg.Name); err != nil {
				log.Println("[registry] failed to report healthy state:", err)
			}
			time.Sleep(1 * time.Second)
		}
	}
}
