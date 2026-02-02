package ratelimit

import (
	"context"
	"time"

	"github.com/syralon/coconut/toolkit/cache"
	"golang.org/x/time/rate"
)

type Option func(*RateLimiter)

func WithInterval(interval time.Duration) Option {
	return func(rl *RateLimiter) {
		if interval > time.Second {
			rl.interval = interval
		}
	}
}

func WithLimit(limit int) Option {
	return func(rl *RateLimiter) {
		if limit > 0 {
			rl.limit = limit
		}
	}
}

func WithSize(size int) Option {
	return func(rl *RateLimiter) {
		if size > 0 {
			rl.size = size
		}
	}
}

type RateLimiter struct {
	cache    *cache.LRUCacheT[string, *rate.Limiter]
	interval time.Duration
	limit    int
	size     int
}

func NewRateLimiter(opts ...Option) *RateLimiter {
	rl := &RateLimiter{
		interval: time.Second,
		limit:    10,
		size:     4096,
	}
	for _, opt := range opts {
		opt(rl)
	}
	rl.cache = cache.NewLRUCacheT[string, *rate.Limiter](cache.WithMaxAge(rl.interval+time.Second), cache.WithMaxSize(rl.size))
	return rl
}

func (rl *RateLimiter) get(key string) *rate.Limiter {
	limiter, ok := rl.cache.Get(key)
	if !ok {
		limiter = rate.NewLimiter(rate.Every(rl.interval), rl.size)
		rl.cache.Put(key, limiter)
	}
	return limiter
}

func (rl *RateLimiter) Allow(key string) bool {
	return rl.get(key).Allow()
}

func (rl *RateLimiter) Wait(ctx context.Context, key string) error {
	return rl.get(key).Wait(ctx)
}
