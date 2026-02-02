package middleware

import (
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/syralon/coconut/internal/ratelimit"
)

func ConcurrencyLimitMiddleware(limit int, match func(r *http.Request) bool, eh ...func(w http.ResponseWriter, r *http.Request, err error)) runtime.Middleware {
	limiter := ratelimit.NewConcurrencyLimiter(limit)
	return func(next runtime.HandlerFunc) runtime.HandlerFunc {
		return func(writer http.ResponseWriter, request *http.Request, params map[string]string) {
			if match(request) {
				if err := limiter.Acquire(request.Context()); err != nil {
					onError(writer, request, err, eh...)
					return
				}
				defer limiter.Release()
			}
			next(writer, request, params)
		}
	}
}

func RateLimitMiddleware(interval time.Duration, limit int, match func(r *http.Request) (string, bool), eh ...func(w http.ResponseWriter, r *http.Request, err error)) runtime.Middleware {
	limiter := ratelimit.NewRateLimiter(ratelimit.WithInterval(interval), ratelimit.WithLimit(limit))
	return func(next runtime.HandlerFunc) runtime.HandlerFunc {
		return func(writer http.ResponseWriter, request *http.Request, params map[string]string) {
			key, ok := match(request)
			if ok {
				if err := limiter.Wait(request.Context(), key); err != nil {
					onError(writer, request, err, eh...)
					return
				}
			}
			next(writer, request, params)
		}
	}
}

func onError(writer http.ResponseWriter, request *http.Request, err error, eh ...func(w http.ResponseWriter, r *http.Request, err error)) {
	if len(eh) > 0 {
		eh[0](writer, request, err)
		return
	}
	http.Error(writer, "To Many Request", http.StatusTooManyRequests)
}
