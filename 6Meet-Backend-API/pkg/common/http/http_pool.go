package http

import (
	"context"
	"net/http"
	"sync"
	"time"
)

type HTTPClientPool struct {
	client *http.Client
	mu     sync.RWMutex
	cache  map[string]any
}

type HTTPClientConfig struct {
	Timeout         time.Duration
	MaxIdleConns    int
	IdleConnTimeout time.Duration
	MaxConnsPerHost int
	EnableCache     bool
	CacheExpiration time.Duration
}

// DefaultHTTPConfig returns default configuration for HTTP client pool
func DefaultHTTPConfig() *HTTPClientConfig {
	return &HTTPClientConfig{
		Timeout:         30 * time.Second,
		MaxIdleConns:    100,
		IdleConnTimeout: 90 * time.Second,
		MaxConnsPerHost: 10,
		EnableCache:     true,
		CacheExpiration: 5 * time.Minute,
	}
}

// NewHTTPClientPool creates a new HTTP client pool with the given configuration
func NewHTTPClientPool(config *HTTPClientConfig) *HTTPClientPool {
	if config == nil {
		config = DefaultHTTPConfig()
	}

	client := &http.Client{
		Timeout: config.Timeout,
		Transport: &http.Transport{
			MaxIdleConns:        config.MaxIdleConns,
			MaxIdleConnsPerHost: config.MaxConnsPerHost,
			IdleConnTimeout:     config.IdleConnTimeout,
		},
	}

	return &HTTPClientPool{
		client: client,
		cache:  make(map[string]any),
	}
}

// RequestWithRetry performs an HTTP request with retry logic
func (p *HTTPClientPool) RequestWithRetry(ctx context.Context, req *http.Request, maxRetries int) (*http.Response, error) {
	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			resp, err := p.client.Do(req)
			if err == nil && resp.StatusCode < 500 {
				return resp, nil
			}
			if err != nil {
				lastErr = err
			}
			// Exponential backoff
			backoff := time.Duration(1<<attempt) * time.Second
			timer := time.NewTimer(backoff)
			select {
			case <-ctx.Done():
				timer.Stop()
				return nil, ctx.Err()
			case <-timer.C:
				continue
			}
		}
	}
	return nil, lastErr
}

// GetFromCache retrieves data from cache if available
func (p *HTTPClientPool) GetFromCache(key string) (any, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	val, ok := p.cache[key]
	return val, ok
}

// SetCache stores data in cache
func (p *HTTPClientPool) SetCache(key string, value any) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.cache[key] = value
}

// ClearCache removes all items from cache
func (p *HTTPClientPool) ClearCache() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.cache = make(map[string]any)
}
