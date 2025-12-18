package coconut

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/syralon/coconut/transport"
)

type Option func(app *App)

func WithReleaser(releasers ...Releaser) Option {
	return func(app *App) {
		app.releasers = append(app.releasers, releasers...)
	}
}

func WithHooks(hooks ...transport.ServerHook) Option {
	return func(app *App) {
		app.hooks = append(app.hooks, hooks...)
	}
}

type App struct {
	awaiting  time.Duration
	releasers Releasers
	hooks     []transport.ServerHook

	servers []transport.Server
}

func NewApp(options ...Option) *App {
	app := &App{awaiting: time.Second}
	for _, option := range options {
		option(app)
	}
	return app
}

func (a *App) Add(server transport.Server, hooks ...transport.ServerHook) *App {
	server = transport.WithHooks(server, append(a.hooks, hooks...)...)
	a.servers = append(a.servers, server)
	a.releasers = append(a.releasers, server)
	return a
}

func (a *App) Run(ctx context.Context) (err error) {
	defer func() { err = errors.Join(err, a.shutdown(ctx)) }()
	done := make(chan struct{}, 1)
	defer close(done)
	if err = a.serves(ctx); err != nil {
		return err
	}
	a.waiting()
	return nil
}

func (a *App) shutdown(ctx context.Context) error {
	timeout := 5 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return a.releasers.Shutdown(ctx)
}

func (a *App) serves(ctx context.Context) error {
	ech := make(chan error, 1)

	for _, server := range a.servers {
		s := server
		go func() {
			err := s.Serve(ctx)
			if err != nil {
				slog.ErrorContext(ctx, err.Error())
				select {
				case ech <- err:
				default:
				}
			}
		}()
	}

	select {
	case err := <-ech:
		return err
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(a.awaiting):
		return nil
	}
}

func (a *App) waiting() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(ch)
	<-ch
}
