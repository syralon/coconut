package logs

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/trace"
)

type otelHandler struct {
	slog.Handler
}

func (h *otelHandler) Handle(ctx context.Context, record slog.Record) error {
	span := trace.SpanContextFromContext(ctx)
	if span.HasSpanID() {
		record.AddAttrs(slog.String("span_id", span.SpanID().String()))
	}
	if span.HasTraceID() {
		record.AddAttrs(slog.String("trace_id", span.TraceID().String()))
	}
	return h.Handler.Handle(ctx, record)
}

func OTEL() LogHook {
	return func(next slog.Handler) slog.Handler {
		return &otelHandler{next}
	}
}
