package initialize

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/global"
	driverHttp "github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/adapters/driver/http"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/common/http/middlewares"
)

// RouterGroup contains all routes
type RouterGroup struct {
	UserHandler driverHttp.UserHandler
}

// NewRouterGroup creates a new RouterGroup
func NewRouterGroup(
	userHandler driverHttp.UserHandler,
) *RouterGroup {
	return &RouterGroup{
		UserHandler: userHandler,
	}
}

// registerRoutes registers all routes
func (rg *RouterGroup) registerRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")

	// User routes
	users := api.Group("/users")
	{
		users.POST("/search", rg.UserHandler.Find)
		users.GET("/:id", rg.UserHandler.Get)

		users.POST("", rg.UserHandler.Create)
		users.PUT("/:id", rg.UserHandler.Update)
		users.DELETE("/:id", rg.UserHandler.Delete)
	}
}

// Ping
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "I'm running!",
	})
}

// NewEngine creates and configures the Gin engine
func NewEngine(routerGroup *RouterGroup) *gin.Engine {
	if global.Config.Server.Mode != "release" {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.New()

	// middlewares
	r.Use(middlewares.CORSMiddleware)

	r.GET("/ping", Ping)

	// Register routes
	routerGroup.registerRoutes(r)

	return r
}
