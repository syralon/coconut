package grpc

import (
	stdgrpc "google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/syralon/coconut/example/internal/config"
	"github.com/syralon/coconut/transport"
	"github.com/syralon/coconut/transport/grpc"
)

func New(c *config.Config) transport.Server {
	srv := grpc.NewServer(&c.GRPC)
	srv.Register(func(srv *stdgrpc.Server) {
		grpc_health_v1.RegisterHealthServer(srv, health.NewServer())
	})
	srv.WithOTELHandler()
	return srv
}
