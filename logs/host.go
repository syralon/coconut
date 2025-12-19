package logs

import (
	"context"
	"log/slog"
	"os"
)

type hostInfoHook struct {
	slog.Handler
	hostname string
}

func (h *hostInfoHook) Handle(ctx context.Context, record slog.Record) error {
	record.AddAttrs(slog.String("hostname", h.hostname))
	return h.Handler.Handle(ctx, record)
}

func HostInfo() LogHook {
	return func(next slog.Handler) slog.Handler {
		hostname, _ := os.Hostname()
		return &hostInfoHook{Handler: next, hostname: hostname}
	}
}
