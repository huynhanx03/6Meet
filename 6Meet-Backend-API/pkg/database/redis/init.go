package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	redisV9 "github.com/redis/go-redis/v9"

	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/settings"
)

const (
	_maxRetries      = 5
	_minRetryBackoff = 300 * time.Millisecond
	_maxRetryBackoff = 500 * time.Millisecond
	_dialTimeout     = 5 * time.Second
	_readTimeout     = 5 * time.Second
	_writeTimeout    = 5 * time.Second
	_minIdleConns    = 20
	_poolTimeout     = 6 * time.Second
	_poolSize        = 300
	_database        = 0
)

type RedisEngine struct {
	client  *redisV9.Client
	config  *settings.RedisConfig
	rwMutex sync.Mutex
}

var ctx = context.Background()

// connect initializes the Redis client
func (r *RedisEngine) connect() error {
	// Build address
	addr := r.config.Host
	if r.config.Port > 0 {
		addr = fmt.Sprintf("%s:%d", addr, r.config.Port)
	}

	// Create redis client
	r.client = redisV9.NewClient(&redisV9.Options{
		Addr:            addr,
		Password:        r.config.Password,
		DB:              r.getOrDefault(r.config.Database, _database),
		PoolSize:        r.getOrDefault(r.config.PoolSize, _poolSize),
		MinIdleConns:    r.getOrDefault(r.config.MinIdleConns, _minIdleConns),
		MaxRetries:      r.getOrDefault(r.config.MaxRetries, _maxRetries),
		DialTimeout:     r.getDurationOrDefault(r.config.DialTimeout, _dialTimeout),
		ReadTimeout:     r.getDurationOrDefault(r.config.ReadTimeout, _readTimeout),
		WriteTimeout:    r.getDurationOrDefault(r.config.WriteTimeout, _writeTimeout),
		MinRetryBackoff: r.getDurationOrDefaultMillis(r.config.MinRetryBackoff, _minRetryBackoff),
		MaxRetryBackoff: r.getDurationOrDefaultMillis(r.config.MaxRetryBackoff, _maxRetryBackoff),
		PoolTimeout:     r.getDurationOrDefault(r.config.PoolTimeout, _poolTimeout),
	})

	// Ping test
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := r.client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis ping failed: %w", err)
	}

	return nil
}

func (r *RedisEngine) getOrDefault(value, def int) int {
	if value > 0 {
		return value
	}
	return def
}

func (r *RedisEngine) getDurationOrDefault(seconds int, def time.Duration) time.Duration {
	if seconds > 0 {
		return time.Duration(seconds) * time.Second
	}
	return def
}

func (r *RedisEngine) getDurationOrDefaultMillis(millis int, def time.Duration) time.Duration {
	if millis > 0 {
		return time.Duration(millis) * time.Millisecond
	}
	return def
}

// Get implements RedisEngine.
func (r *RedisEngine) Get(key string) ([]byte, bool, error) {
	byteValue, err := r.client.Get(ctx, key).Bytes()
	if err == redisV9.Nil {
		return nil, false, err
	}
	if err != nil {
		return nil, false, err
	}
	return byteValue, true, nil
}

// Invalidate implements RedisEngine.
func (r *RedisEngine) Invalidate(key string) error {
	r.rwMutex.Lock()
	defer r.rwMutex.Unlock()
	//Delete Key From Cache
	return r.client.Del(ctx, key).Err()
}

func (r *RedisEngine) InvalidatePrefix(prefix string) error {
	r.rwMutex.Lock()
	defer r.rwMutex.Unlock()
	pattern := prefix + "*"
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}
	for _, key := range keys {
		if err := r.client.Del(ctx, key).Err(); err != nil {
			return err
		}
	}
	return nil
}

// Set implements RedisEngine.
func (r *RedisEngine) Set(key string, value any, timeToLive ...time.Duration) error {
	r.rwMutex.Lock()
	defer r.rwMutex.Unlock()
	byteValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	//set default value 0
	var ttl time.Duration
	if len(timeToLive) > 0 {
		ttl = timeToLive[0]
	}
	//Set value to RedisEngine cache
	return r.client.Set(ctx, key, byteValue, ttl).Err()
}

// Close implements RedisEngine.
func (r *RedisEngine) Close() {
	r.client.Close()
}

// Client implements RedisEngine.
func (r *RedisEngine) Client() *redisV9.Client {
	return r.client
}
