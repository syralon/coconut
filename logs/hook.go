package logs

import "log/slog"

type LogHook func(next slog.Handler) slog.Handler

func WithHooks(h slog.Handler, hooks ...LogHook) slog.Handler {
	for _, hook := range hooks {
		h = hook(h)
	}
	return h
}

func SetDefault(h slog.Handler, hooks ...LogHook) {
	slog.SetDefault(slog.New(WithHooks(h, hooks...)))
}
