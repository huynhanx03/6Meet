
package initialize

import (
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/global"
	db "github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/adapters/driven/db"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/adapters/driver/http"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/core/application/service"
)

// InitializeServer wires up all dependencies and returns the Server
func InitializeServer() *Server {
	// Initialize repositories
	userRepo := db.NewUserRepository(global.MongoDB.DB)

	// Initialize services
	userService := service.NewUserService(userRepo)

	// Initialize controllers
	userHandler := http.NewUserHandler(userService)

	// Create router group with dependencies
	routerGroup := NewRouterGroup(userHandler)

	// Create Gin engine
	engine := NewEngine(routerGroup)

	// Create Server
	return NewServer(engine)
}
