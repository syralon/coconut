package gateway

import (
	"github.com/syralon/coconut/example/internal/config"
	"github.com/syralon/coconut/example/internal/service"
	"github.com/syralon/coconut/example/proto/syralon/example"
	"github.com/syralon/coconut/transport"
	"github.com/syralon/coconut/transport/gateway"
)

func NewServer(c *config.Config, services *service.Services) transport.Server {
	srv := gateway.NewServer(&c.Gateway)
	if c.Gateway.Endpoint != "" {
		srv.RegisterEndpoint(
			c.Gateway.Endpoint,
			example.RegisterBookServiceHandlerFromEndpoint,
			example.RegisterBookShelfServiceHandlerFromEndpoint,
		)
	} else {
		srv.Register(
			gateway.ServerRegister[example.BookServiceServer](services.Book, example.RegisterBookServiceHandlerServer),
			gateway.ServerRegister[example.BookShelfServiceServer](services.BookShelf, example.RegisterBookShelfServiceHandlerServer),
		)
	}
	return srv
}
