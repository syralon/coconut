package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/syralon/coconut/toolkit/cache"
	"github.com/syralon/coconut/toolkit/netutil"
	"golang.org/x/time/rate"
)

type RateLimiter struct {
	cache    *cache.LRUCacheT[string, *rate.Limiter]
	interval time.Duration
	size     int
	match    func(r *http.Request) (string, bool)
}

func NewRateLimiter(interval time.Duration, size int, match func(r *http.Request) (string, bool)) *RateLimiter {
	if interval < time.Second {
		interval = time.Second
	}
	if size < 1 {
		size = 1
	}
	if match == nil {
		match = func(request *http.Request) (string, bool) { return netutil.ClientIP(request), true }
	}
	return &RateLimiter{
		cache:    cache.NewLRUCacheT[string, *rate.Limiter](cache.WithMaxSize(1024), cache.WithMaxAge(time.Minute)),
		interval: interval,
		size:     size,
		match:    match,
	}
}

func (r *RateLimiter) Middleware(mux *runtime.ServeMux) runtime.Middleware {
	return func(next runtime.HandlerFunc) runtime.HandlerFunc {
		return func(writer http.ResponseWriter, request *http.Request, params map[string]string) {
			key, ok := r.match(request)
			if !ok {
				next(writer, request, params)
				return
			}
			if !r.allow(key) {
				_, outboundMarshaler := runtime.MarshalerForRequest(mux, request)
				runtime.HTTPError(request.Context(), mux, outboundMarshaler, writer, request, &runtime.HTTPStatusError{
					HTTPStatus: http.StatusTooManyRequests,
					Err:        context.DeadlineExceeded,
				})
				return
			}
			next(writer, request, params)
		}
	}
}

func (r *RateLimiter) allow(key string) bool {
	limiter, ok := r.cache.Get(key)
	if !ok {
		limiter = rate.NewLimiter(rate.Every(r.interval), r.size)
		r.cache.Put(key, limiter)
	}
	return limiter.Allow()
}

func RateLimitMiddleware(mux *runtime.ServeMux, interval time.Duration, size int, match func(r *http.Request) (string, bool)) runtime.Middleware {
	return NewRateLimiter(interval, size, match).Middleware(mux)
}

func ConcurrencyLimitMiddleware(mux *runtime.ServeMux, concurrency int, match func(r *http.Request) bool) runtime.Middleware {
	if match == nil {
		match = func(*http.Request) bool { return true }
	}
	ch := make(chan struct{}, concurrency)
	return func(next runtime.HandlerFunc) runtime.HandlerFunc {
		return func(writer http.ResponseWriter, request *http.Request, params map[string]string) {
			if match(request) {
				if err := acquire(request.Context(), ch); err != nil {
					_, outboundMarshaler := runtime.MarshalerForRequest(mux, request)
					runtime.HTTPError(request.Context(), mux, outboundMarshaler, writer, request, &runtime.HTTPStatusError{
						HTTPStatus: http.StatusTooManyRequests,
						Err:        err,
					})
					return
				}
				defer release(ch)
			}
			next(writer, request, params)
		}
	}
}

func acquire(ctx context.Context, c chan struct{}) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-c:
		return nil
	}
}

func release(c chan struct{}) { <-c }
