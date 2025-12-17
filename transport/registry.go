package transport

import (
	"context"
	"errors"
	"log/slog"

	"github.com/syralon/coconut/mesh"
)

type serverRegistry struct {
	Server

	registry mesh.Registry
	receipt  *mesh.Receipt
	endpoint *mesh.Endpoint
}

func (s *serverRegistry) Serve(ctx context.Context) (err error) {
	if ed, ok := s.Server.(EndpointServer); !ok {
		return s.Server.Serve(ctx)
	} else {
		s.endpoint = ed.Endpoint()
	}
	if s.receipt, err = s.registry.Register(ctx, s.endpoint); err != nil {
		return err
	}
	slog.InfoContext(ctx, "service registered", "name", s.endpoint.Name, "address", s.endpoint.Address())
	return s.Server.Serve(ctx)
}

func (s *serverRegistry) Shutdown(ctx context.Context) error {
	if s.endpoint == nil {
		return s.Server.Shutdown(ctx)
	}
	err := s.registry.Deregister(ctx, s.endpoint)
	return errors.Join(err, s.Server.Shutdown(ctx))
}

func Registry(registry mesh.Registry) ServerHook {
	return func(server Server) Server {
		return &serverRegistry{
			Server:   server,
			registry: registry,
		}
	}
}
