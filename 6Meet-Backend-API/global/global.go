package global

import (
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/database/mongodb"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/database/redis"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/logger"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/settings"
)

var (
	Version string = "v0.0.1"
	Logger  *logger.LoggerZap
	Config  *settings.Config
	MongoDB *mongodb.Connection
	Redis   *redis.RedisEngine
)
