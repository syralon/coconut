package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	stdgrpc "google.golang.org/grpc"

	"github.com/syralon/coconut"
	"github.com/syralon/coconut/configuration"
	"github.com/syralon/coconut/example/internal/config"
	"github.com/syralon/coconut/logs"
	"github.com/syralon/coconut/transport"
	"github.com/syralon/coconut/transport/gateway"
	"github.com/syralon/coconut/transport/gateway/middleware"
	"github.com/syralon/coconut/transport/grpc"
	"github.com/syralon/coconut/transport/grpc/interceptor"
)

func init() {
	logs.SetDefault(
		slog.NewTextHandler(os.Stdout, nil),
		logs.Gateway(),
		logs.HostInfo(),
	)
}

func main() {
	ctx := context.Background()

	c := new(config.Config)
	if err := configuration.Read(ctx, c); err != nil {
		panic(err)
	}

	app := coconut.NewApp(coconut.WithHooks(transport.Logger()))
	srv1 := gateway.NewServer(&c.Gateway)
	srv1.WithOptions(
		runtime.WithMiddlewares(
			middleware.Logger(),
			middleware.OTEL(),
		),
	)
	srv2 := grpc.NewServer(&c.GRPC)
	srv2.WithUnaryInterceptor(interceptor.RecoveryUnaryServerInterceptor())
	srv2.WithStreamInterceptor(interceptor.RecoveryStreamServerInterceptor())
	srv2.WithOptions(stdgrpc.StatsHandler(interceptor.NewServerHandler()))
	app.Add(srv1).Add(srv2)
	if err := app.Run(ctx); err != nil {
		panic(err)
	}
}
