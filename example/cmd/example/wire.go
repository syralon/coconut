//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"github.com/google/wire"
	"github.com/syralon/coconut"
	"github.com/syralon/coconut/example/internal/api/controller"
	"github.com/syralon/coconut/example/internal/config"
	"github.com/syralon/coconut/example/internal/infrastructure/data"
	"github.com/syralon/coconut/example/internal/infrastructure/dependency"
	"github.com/syralon/coconut/example/internal/server"
)

func initialize(config *config.Config) (*coconut.App, func(), error) {
	panic(wire.Build(
		controller.ProviderSet,
		dependency.ProviderSet,
		server.ProviderSet,
		data.ProviderSet,
		newApp,
	))
}
