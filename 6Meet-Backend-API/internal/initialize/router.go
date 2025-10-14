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
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/container"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/ports"
	middleware "github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/common/http/middlewares"
	"go.uber.org/zap"
)

// RouterGroup contains all routes
type RouterGroup struct {
	UserHandler ports.IUserHandler
}

// NewRouterGroup creates a new RouterGroup
func NewRouterGroup(
	userHandler ports.IUserHandler,
) *RouterGroup {
	return &RouterGroup{
		UserHandler: userHandler,
	}
}

func (rg *RouterGroup) registerRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")

	// User routes
	users := api.Group("/users")
	{
		users.POST("/list", rg.UserHandler.ListUsers())
		users.GET("/:id", rg.UserHandler.GetUserByID())

		users.POST("", rg.UserHandler.CreateUser())
		users.PUT("/:id", rg.UserHandler.UpdateUser())
		users.DELETE("/:id", rg.UserHandler.DeleteUser())
	}
}

// Ping
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "6Meet-Backend-API is running",
	})
}

func InitRouter() {
	if global.Config.Server.Mode != "release" {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.New()

	// middlewares
	r.Use(middleware.CORSMiddleware)

	r.GET("/ping", Ping)

	// Register routes
	routerGroup := NewRouterGroup(
		container.GetDependencies().UserHandler,
	)
	routerGroup.registerRoutes(r)

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
		Handler: r,
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
	}
}
