package http

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/alexraileanu/thomas-appie/pkg/db"
)

type Server struct {
	engine    *echo.Echo
	dbService *db.Service
}

func NewServer(dbService *db.Service) *Server {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))
	server := &Server{
		engine:    e,
		dbService: dbService,
	}

	server.registerRoutes()

	return server
}

func (s *Server) Start() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	go func() {
		if err := s.engine.Start(":7008"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.engine.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.engine.Shutdown(ctx); err != nil {
		s.engine.Logger.Fatal(err)
	}
}

func (s *Server) registerRoutes() {
	apiGroup := s.engine.Group("/api")
	apiGroup.GET("/products", s.containersListHandler)
}
