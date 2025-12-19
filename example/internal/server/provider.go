package server

import (
	"github.com/google/wire"
	"github.com/syralon/coconut/example/internal/api/controller"
	"github.com/syralon/coconut/example/internal/config"
	"github.com/syralon/coconut/example/internal/server/gateway"
	"github.com/syralon/coconut/example/internal/server/grpc"
	"github.com/syralon/coconut/transport"
)

type Servers []transport.Server

func NewServers(c *config.Config, services *controller.Services) Servers {
	return Servers{
		gateway.NewServer(c, services),
		grpc.NewServer(c, services),
	}
}

var ProviderSet = wire.NewSet(
	NewServers,
)
