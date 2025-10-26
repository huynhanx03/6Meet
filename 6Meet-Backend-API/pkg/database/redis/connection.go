package redis

import (
	"fmt"

	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/settings"
)

// NewConnection creates and returns a new Redis client
func NewConnection(cfg *settings.RedisConfig) (*RedisEngine, error) {
	engine := &RedisEngine{
		config: cfg,
	}

	if err := engine.connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return engine, nil
}
