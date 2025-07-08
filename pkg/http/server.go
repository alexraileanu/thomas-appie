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

	"github.com/alexraileanu/thomas-appie/pkg/config"
	"github.com/alexraileanu/thomas-appie/pkg/db"
	"github.com/alexraileanu/thomas-appie/pkg/logger"
	"github.com/alexraileanu/thomas-appie/pkg/product"
	"github.com/alexraileanu/thomas-appie/web"
)

type Server struct {
	engine         *echo.Echo
	dbService      *db.Service
	productService *product.Service
	loggerService  *logger.Service

	conf config.Appie
}

func NewServer(dbService *db.Service, productService *product.Service, conf config.Appie, loggerService *logger.Service) *Server {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))
	server := &Server{
		engine:         e,
		dbService:      dbService,
		productService: productService,
		loggerService:  loggerService,

		conf: conf,
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

	apiGroup.POST("/products/refresh", s.refreshProducts)
	apiGroup.GET("/products", s.getDiscountedProducts)

	apiGroup.GET("/db/products", s.getAllProducts)
	apiGroup.POST("/db/products", s.updateProducts)

	assetHandler := http.FileServer(web.Dist())
	s.engine.GET("/", echo.WrapHandler(assetHandler))
	s.engine.GET("/assets/*", echo.WrapHandler(assetHandler))
}
