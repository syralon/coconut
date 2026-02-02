package interceptor

import (
	"context"
	"time"

	"github.com/syralon/coconut/internal/ratelimit"
	"google.golang.org/grpc"
)

func ConcurrencyLimitUnaryInterceptor(limit int, match func(ctx context.Context) bool) grpc.UnaryServerInterceptor {
	limiter := ratelimit.NewConcurrencyLimiter(limit)
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (rsp any, err error) {
		if match(ctx) {
			if err = limiter.Acquire(ctx); err != nil {
				return nil, err
			}
			defer limiter.Release()
		}
		return handler(ctx, req)
	}
}

func RateLimitUnaryInterceptor(interval time.Duration, limit int, match func(ctx context.Context) (string, bool)) grpc.UnaryServerInterceptor {
	limiter := ratelimit.NewRateLimiter(ratelimit.WithInterval(interval), ratelimit.WithLimit(limit))
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		if key, ok := match(ctx); ok {
			if err = limiter.Wait(ctx, key); err != nil {
				return nil, err
			}
		}
		return handler(ctx, req)
	}
}
