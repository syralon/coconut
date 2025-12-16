package logs

import (
	"context"
	"log/slog"
	"os"

	"github.com/syralon/coconut/toolkit/netutil"
)

type hostInfoHook struct {
	slog.Handler
	hostname string
	ip       string
}

func (h *hostInfoHook) Handle(ctx context.Context, record slog.Record) error {
	record.AddAttrs(slog.String("hostname", h.hostname), slog.String("ip", h.ip))
	return h.Handler.Handle(ctx, record)
}

func HostInfo() LogHook {
	return func(next slog.Handler) slog.Handler {
		hostname, _ := os.Hostname()
		ip, _ := netutil.InternalIPV4()
		return &hostInfoHook{Handler: next, hostname: hostname, ip: ip}
	}
}
