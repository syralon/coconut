package transport

import (
	"context"

	"github.com/syralon/coconut/mesh"
)

type Server interface {
	Serve(ctx context.Context) error
	Shutdown(ctx context.Context) error
	Endpoint() *mesh.Endpoint
}

type ServerHook func(server Server) Server

func WithHooks(server Server, hooks ...ServerHook) Server {
	for _, h := range hooks {
		server = h(server)
	}
	return server
}
