package etcdutil

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"

	"github.com/syralon/coconut/mesh"
)

type ETCDElector struct {
	id     string
	client *clientv3.Client
	opts   []concurrency.SessionOption
}

func (e *ETCDElector) Start(ctx context.Context, name string, runner mesh.ElectionRunner) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := e.campaign(ctx, name, runner); err != nil {
				slog.ErrorContext(ctx, err.Error())
				time.Sleep(time.Second)
			}
		}
	}
}

func (e *ETCDElector) Close() {
}

func (e *ETCDElector) campaign(ctx context.Context, name string, runner mesh.ElectionRunner) error {
	session, err := concurrency.NewSession(e.client)
	if err != nil {
		return err
	}
	defer func() {
		if err = session.Close(); err != nil {
			slog.ErrorContext(ctx, err.Error())
		}
	}()
	election := concurrency.NewElection(session, name)
	if err = election.Campaign(ctx, e.id); err != nil {
		return err
	}
	defer func() {
		err = election.Resign(ctx)
		if err != nil {
			slog.ErrorContext(ctx, err.Error())
		}
	}()
	defer runner.Close()
	if err = runner.Run(ctx); err != nil {
		slog.ErrorContext(ctx, err.Error())
	}
	return nil
}

func NewETCDElector(client *clientv3.Client, options ...concurrency.SessionOption) *ETCDElector {
	return &ETCDElector{
		id:     uuid.New().String(),
		client: client,
		opts:   options,
	}
}
