package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/syralon/coconut"
	"github.com/syralon/coconut/configuration"
	"github.com/syralon/coconut/example/internal/config"
	"github.com/syralon/coconut/example/internal/global/snowflake"
	"github.com/syralon/coconut/example/internal/server"
	"github.com/syralon/coconut/logs"
	"github.com/syralon/coconut/pkg/etcdutil"
	"github.com/syralon/coconut/transport"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func init() {
	logs.SetDefault(
		slog.NewTextHandler(os.Stdout, nil),
		logs.Gateway(),
		logs.HostInfo(),
	)
}

func newApp(client *clientv3.Client, servers server.Servers) (*coconut.App, error) {
	ctx := context.Background()

	roulette := etcdutil.NewRoulette("example", client)
	id, err := roulette.Allocate(ctx)
	if err != nil {
		return nil, err
	}
	if err = snowflake.Setup(id); err != nil {
		return nil, err
	}
	app := coconut.NewApp(
		coconut.WithHooks(
			transport.Logger(),
			transport.Registry(etcdutil.NewRegistry(client)),
		),
		coconut.WithReleaser(
			coconut.ReleaserFunc(roulette.Release),
		),
	)
	app.Add(servers...)
	return app, nil
}

func main() {
	ctx := context.Background()

	c := new(config.Config)
	if err := configuration.Read(ctx, c); err != nil {
		panic(err)
	}

	//client, err := c.ETCD.NewClient()
	//if err != nil {
	//	panic(err)
	//}
	//
	//app := coconut.NewApp(
	//	coconut.WithHooks(
	//		transport.Logger(),
	//		transport.Registry(etcdutil.NewRegistry(client)),
	//	),
	//	coconut.WithReleaser(
	//		coconut.ReleaserFunc(roulette.Release),
	//	),
	//)
	//
	//srv1 := gateway.NewServer(&c.Gateway)
	//srv1.WithOptions(
	//	runtime.WithMiddlewares(
	//		middleware.Logger(),
	//		middleware.OTEL(),
	//	),
	//)
	//srv2 := grpc.NewServer(&c.GRPC)
	//srv2.WithUnaryInterceptor(interceptor.RecoveryUnaryServerInterceptor())
	//srv2.WithStreamInterceptor(interceptor.RecoveryStreamServerInterceptor())
	//srv2.WithOptions(stdgrpc.StatsHandler(interceptor.NewServerHandler()))
	//app.Add(srv1).Add(srv2)
	//if err := app.Run(ctx); err != nil {
	//	panic(err)
	//}
}
