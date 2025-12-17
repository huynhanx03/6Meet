package initialize

import (
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/global"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/database/redis"
)

// SetupRedis initializes the Redis connection
func SetupRedis() {
    config := global.Config.Redis

    engine, err := redis.NewConnection(&config)

    if err != nil {
        global.Logger.Sugar().Fatalf("Failed to connect to Redis: %v", err)
    }

    global.Redis = engine
    global.Logger.Sugar().Info("Connected to Redis successfully")
}