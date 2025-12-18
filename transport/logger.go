package transport

import (
	"context"
	"fmt"
	"log/slog"
)

type serverLogger struct {
	Server
}

func (s *serverLogger) Serve(ctx context.Context) (err error) {
	slog.InfoContext(ctx, fmt.Sprintf("%s server started", s.Name()))
	return s.Server.Serve(ctx)
}

func (s *serverLogger) Shutdown(ctx context.Context) error {
	slog.InfoContext(ctx, fmt.Sprintf("%s server stopped", s.Name()))
	return s.Server.Shutdown(ctx)
}

func Logger() ServerHook {
	return func(server Server) Server {
		return &serverLogger{
			Server: server,
		}
	}
}
