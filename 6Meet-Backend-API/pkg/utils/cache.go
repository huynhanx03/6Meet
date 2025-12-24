package utils

import (
	"context"
	"encoding/json"

	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/common/cache"
	"github.com/pkg/errors"
)

func HandleHitCache(ctx context.Context, model any, c cache.CacheEngine, key string) error {
	byteData, exists, err := c.Get(ctx, key)
	if exists && err == nil {
		err = json.Unmarshal(byteData, model)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal cache")
		}
		return nil
	}
	return errors.Wrap(err, "miss cache")
}

func HandleSetCache(ctx context.Context, model any, c cache.CacheEngine, key string, ttl int) error {
	return c.Set(ctx, key, model, ToDuration(ttl))
}
