package initialize

import (
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/global"
	driven "github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/adapters/driven/db"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/adapters/driver/http"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/container"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/core/service"
)

// InitDependencies initializes all dependencies
func InitDependencies() {
	// Initialize repositories
	userRepo := driven.NewUserRepository(global.MongoDB.DB)

	// Initialize services
	userService := service.NewUserService(userRepo)

	// Initialize controllers
	userHandler := http.NewUserHandler(userService)

	deps := &container.DependencyContainer{
		// Repositories
		UserRepo:       userRepo,

		// Services
		UserService:    userService,

		// Controllers
		UserHandler:    userHandler,
	}

	// Set as singleton
	container.SetDependencies(deps)
}
