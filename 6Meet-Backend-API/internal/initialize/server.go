package initialize

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/global"
	"go.uber.org/zap"
)

// Server wraps the HTTP server
type Server struct {
	engine *gin.Engine
}

// NewServer creates a new Server instance
func NewServer(engine *gin.Engine) *Server {
	return &Server{
		engine: engine,
	}
}

// Run starts the HTTP server with graceful shutdown
func (s *Server) Run() error {
	host := os.Getenv("SERVER_HOST")
	if host == "" {
		host = global.Config.Server.Host
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = strconv.Itoa(global.Config.Server.Port)
	}

	address := fmt.Sprintf("%s:%s", host, port)

	srv := &http.Server{
		Addr:    address,
		Handler: s.engine,
	}

	// Start server in a goroutine
	go func() {
		global.Logger.Info("Server starting",
			zap.String("address", address),
			zap.String("mode", global.Config.Server.Mode),
		)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			global.Logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	global.Logger.Info("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		global.Logger.Error("Server forced to shutdown", zap.Error(err))
		return err
	}
	
	global.Logger.Info("Server exited")
	return nil
}
