package middleware

import (
	"log/slog"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/syralon/coconut/proto/syralon/coconut/errors"
)

func Recovery() runtime.Middleware {
	return func(next runtime.HandlerFunc) runtime.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			defer func() {
				if err := errors.Recovery(recover()); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					slog.ErrorContext(r.Context(), err.Error(), slog.String("stack", err.Stack()))

					return
				}
			}()
			next(w, r, pathParams)
		}
	}
}
