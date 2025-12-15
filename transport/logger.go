package transport

import (
	"context"
	"log/slog"

	"github.com/syralon/coconut/mesh"
)

type serverLogger struct {
	Server
	endpoint *mesh.Endpoint
}

func (s *serverLogger) Serve(ctx context.Context) (err error) {
	s.endpoint = s.Endpoint()
	slog.InfoContext(ctx, "service started", "name", s.endpoint.Name, "endpoint", s.endpoint.Address())
	return s.Server.Serve(ctx)
}

func (s *serverLogger) Shutdown(ctx context.Context) error {
	defer slog.InfoContext(ctx, "service stopped", "name", s.endpoint.Name, "endpoint", s.endpoint.Address())
	return s.Server.Shutdown(ctx)
}

func Logger() ServerHook {
	return func(server Server) Server {
		return &serverLogger{
			Server: server,
		}
	}
}
