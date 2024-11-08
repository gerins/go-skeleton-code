package cmd

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"go-skeleton-code/config"
	"go-skeleton-code/pkg/log"
	middlewareLog "go-skeleton-code/pkg/log/middleware/gin"
)

type HTTPServer struct {
	Server *gin.Engine
	cfg    *config.Config
}

// NewHTTPServer initializes the Gin HTTP server with configuration.
func NewHTTPServer(cfg *config.Config) *HTTPServer {
	// Optional: set Gin to release mode for production
	gin.SetMode(gin.ReleaseMode)
	server := gin.New()

	// Apply middleware
	server.Use(gin.Recovery())                 // Built-in Gin recovery middleware
	server.Use(middlewareLog.SetLogRequest())  // Custom Logging middleware
	server.Use(middlewareLog.SaveLogRequest()) // Custom Request and response logging

	return &HTTPServer{
		cfg:    cfg,
		Server: server,
	}
}

func (s *HTTPServer) Run() chan bool {
	// Health check route
	s.Server.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Start server
	var (
		serverExitSignal = make(chan bool)
		address          = fmt.Sprintf("%v:%v", s.cfg.App.HTTP.Host, s.cfg.App.HTTP.Port)
		httpServer       = &http.Server{
			Addr:    address,
			Handler: s.Server.Handler(),
		}
	)

	log.Infof("%v server app and running %v", s.cfg.App.Name, address)

	go func() {
		// service connections
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("%v server, %v", s.cfg.App.Name, err)
		}
	}()

	// Graceful shutdown on interrupt
	go func() {
		<-serverExitSignal
		log.Info("stopping http server")

		// Wait for server to complete outstanding requests
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(ctx); err != nil {
			log.Fatalf("failed stopping server, %v", err)
		}

		log.Info("finished stopping http server")
		serverExitSignal <- true
	}()

	return serverExitSignal
}
