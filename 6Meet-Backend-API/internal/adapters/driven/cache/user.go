package cache

import (
	"context"

	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/constant"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/core/dto"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/core/entity"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/core/mapper"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/ports"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/common/cache"
)

type userCache struct {
	redis cache.CacheEngine
}

// NewUser creates a new user cache repository
func NewUser(redis cache.CacheEngine) ports.UserCacheRepository {
	return &userCache{
		redis: redis,
	}
}

// getKey returns the key for a user in the cache
func (c *userCache) getKey(id string) string {
	return constant.UserCachePrefix + id
}

// Set stores a user in the cache
func (c *userCache) Set(ctx context.Context, user *entity.User) error {
	return cache.HandleSetCache(ctx, user, c.redis, c.getKey(user.ID), constant.UserCacheTTL)
}

// Get retrieves a user from the cache
func (c *userCache) Get(ctx context.Context, id string) (*dto.UserResponse, error) {
	var user entity.User

	if err := cache.HandleHitCache(ctx, &user, c.redis, c.getKey(id)); err != nil {
		return nil, err
	}

	return mapper.ToUserResponse(&user), nil
}

// BatchSet stores multiple users in the cache using pipeline
func (c *userCache) BatchSet(ctx context.Context, users []*entity.User) error {
	values := make(map[string]any, len(users))
	for _, user := range users {
		values[c.getKey(user.ID)] = user
	}
	return c.redis.BatchSet(ctx, values, constant.UserCacheTTL)
}

// BatchDelete removes multiple users from the cache
func (c *userCache) BatchDelete(ctx context.Context, ids []string) error {
	keys := make([]string, len(ids))
	for i, id := range ids {
		keys[i] = c.getKey(id)
	}
	return c.redis.BatchDelete(ctx, keys)
}

func (c *userCache) Delete(ctx context.Context, id string) error {
	return c.redis.Delete(ctx, c.getKey(id))
}
