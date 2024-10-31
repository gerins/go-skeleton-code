package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	return &HTTPServer{
		cfg:    cfg,
		Server: gin.New(),
	}
}

func (s *HTTPServer) Run() chan bool {
	// Apply middleware
	s.Server.Use(gin.Recovery())                 // Built-in Gin recovery middleware
	s.Server.Use(middlewareLog.SetLogRequest())  // Custom Logging middleware
	s.Server.Use(middlewareLog.SaveLogRequest()) // Custom Request and response logging

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

	go func() {
		// service connections
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("%v server, %v", s.cfg.App.Name, err)
		}

		log.Infof("%v server app and running, %v", s.cfg.App.Name, address)
	}()

	// Graceful shutdown on interrupt
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

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
