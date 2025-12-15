package logs

import (
	"context"
	"log/slog"
)

type kvHandler struct {
	slog.Handler
	kv map[string]string
}

func (h *kvHandler) Handle(ctx context.Context, record slog.Record) error {
	for k, v := range h.kv {
		record.AddAttrs(slog.String(k, v))
	}
	return h.Handler.Handle(ctx, record)
}

func KVParis(kv map[string]string) LogHook {
	return func(next slog.Handler) slog.Handler {
		return &kvHandler{Handler: next, kv: kv}
	}
}
