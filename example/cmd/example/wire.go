//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"github.com/google/wire"
	"github.com/syralon/coconut"
	"github.com/syralon/coconut/example/internal/config"
	"github.com/syralon/coconut/example/internal/infra"
	"github.com/syralon/coconut/example/internal/server"
	"github.com/syralon/coconut/example/internal/service"
)

func initialize(config *config.Config) (*coconut.App, func(), error) {
	panic(wire.Build(
		infra.ProviderSet,
		service.ProviderSet,
		server.ProviderSet,
		newApp,
	))
}
