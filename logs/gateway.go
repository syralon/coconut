package logs

import (
	"context"
	"log/slog"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type gatewayHandler struct {
	slog.Handler
}

func (h *gatewayHandler) Handle(ctx context.Context, record slog.Record) error {
	if pattern, ok := runtime.HTTPPathPattern(ctx); ok {
		record.AddAttrs(slog.String("pattern", pattern))
	}
	if method, ok := runtime.RPCMethod(ctx); ok {
		record.AddAttrs(slog.String("method", method))
	}
	return h.Handler.Handle(ctx, record)
}

func Gateway() LogHook {
	return func(next slog.Handler) slog.Handler {
		return &gatewayHandler{next}
	}
}
