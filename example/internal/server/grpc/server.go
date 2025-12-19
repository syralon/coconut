package grpc

import (
	"github.com/syralon/coconut/example/internal/api/controller"
	"github.com/syralon/coconut/example/internal/config"
	"github.com/syralon/coconut/example/proto/syralon/example"
	"github.com/syralon/coconut/transport"
	"github.com/syralon/coconut/transport/grpc"
	stdgrpc "google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func NewServer(c *config.Config, services *controller.Services) transport.Server {
	srv := grpc.NewServer(&c.GRPC)
	srv.Register(func(srv *stdgrpc.Server) {
		grpc_health_v1.RegisterHealthServer(srv, health.NewServer())
		example.RegisterBookServiceServer(srv, services.BookService)
		example.RegisterBookShelfServiceServer(srv, services.BookShelfService)
	})
	srv.WithOTELHandler()
	return srv
}
