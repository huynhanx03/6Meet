package container

import (
	"sync"

	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/ports"
)

// DependencyContainer holds all dependencies
type DependencyContainer struct {
	// Repositories
	UserRepo ports.IUserRepository

	// Services
	UserService ports.IUserService

	// Handlers
	UserHandler ports.IUserHandler
}

var (
	depsContainer *DependencyContainer
	depsInitOnce  sync.Once
)

// GetDependencies returns the singleton instance of DependencyContainer
func GetDependencies() *DependencyContainer {
	return depsContainer
}

// SetDependencies sets the singleton instance (should only be called once during initialization)
func SetDependencies(deps *DependencyContainer) {
	depsInitOnce.Do(func() {
		depsContainer = deps
	})
}

// IsDependenciesInitialized checks if the singleton has been initialized
func IsDependenciesInitialized() bool {
	return depsContainer != nil
}
