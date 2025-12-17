package mesh

import (
	"context"
)

type ElectionRunner interface {
	Run(ctx context.Context) error
	Close()
}

type Elector interface {
	Start(ctx context.Context, name string, runner ElectionRunner) error
	Close()
}

type ElectionServer struct {
	name    string
	elector Elector
	runner  ElectionRunner
}

func NewElectionServer(name string, elector Elector, runner ElectionRunner) *ElectionServer {
	return &ElectionServer{name: name, elector: elector, runner: runner}
}

func (e *ElectionServer) Serve(ctx context.Context) error {
	return e.elector.Start(ctx, e.name, e.runner)
}

func (e *ElectionServer) Shutdown(_ context.Context) error {
	e.elector.Close()
	return nil
}
