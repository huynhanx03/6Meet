package di

import (
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/adapters/driver/http"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/ports"
)

// Container holds all dependencies for the application
type Container struct {
	UserRepo ports.UserRepository
	UserSearch      ports.UserSearchRepository
	UserCache           ports.UserCacheRepository
	UserService         ports.UserService
	UserHandler         http.UserHandler
	CrawlService        ports.CrawlService
	UserCrawlerConsumer ports.UserConsumer
	UserCDCConsumer     ports.UserConsumer
	UserProducer        ports.UserProducer
}

// GlobalContainer is the global instance of the dependency container
var GlobalContainer *Container
