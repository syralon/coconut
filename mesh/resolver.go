package mesh

import (
	"context"
	"log/slog"
	"strings"
	"sync"

	gresolver "google.golang.org/grpc/resolver"
)

type BuilderOption func(*builder)

func WithBalancer(balancer Balancer) BuilderOption {
	return func(b *builder) {
		b.balancer = balancer
	}
}

func NewBuilder(scheme string, discovery Discovery, opts ...BuilderOption) gresolver.Builder {
	b := &builder{
		scheme:    scheme,
		discovery: discovery,
	}
	for _, opt := range opts {
		opt(b)
	}
	return b
}

type builder struct {
	discovery Discovery
	balancer  Balancer
	scheme    string
}

func (b builder) Build(target gresolver.Target, cc gresolver.ClientConn, _ gresolver.BuildOptions) (gresolver.Resolver, error) {
	// Refer to https://github.com/grpc/grpc-go/blob/16d3df80f029f57cff5458f1d6da6aedbc23545d/clientconn.go#L1587-L1611
	// grpc://etcd.example.com/example_service
	// grpc:///example_service
	endpoint := target.URL.Path
	if endpoint == "" {
		endpoint = target.URL.Opaque
	}
	endpoint = strings.TrimPrefix(endpoint, "/")
	r := &resolver{
		discovery: b.discovery,
		target:    endpoint,
		cc:        cc,
		balancer:  b.balancer,
		done:      make(chan struct{}),
	}
	ctx := context.Background()
	ctx, r.cancel = context.WithCancel(ctx)

	return r, r.resolve(ctx)
}

func (b builder) Scheme() string {
	return b.scheme
}

type resolver struct {
	discovery Discovery
	target    string
	cc        gresolver.ClientConn
	balancer  Balancer

	done   chan struct{}
	cancel context.CancelFunc
	ups    sync.Map
}

func (r *resolver) resolve(ctx context.Context) error {
	endpoints, err := r.discovery.Discover(ctx, r.target)
	if err != nil {
		return err
	}

	if r.balancer != nil {
		endpoints = endpoints.Balance(ctx, r.balancer)
	}

	for _, endpoint := range endpoints {
		if endpoint.State != StateUp {
			continue
		}
		addr := endpoint.Address()
		if addr == "" {
			continue
		}
		r.ups.Store(addr, struct{}{})
	}

	if err = r.update(); err != nil {
		return err
	}
	go r.watch(ctx)
	return nil
}

func (r *resolver) watch(ctx context.Context) {
	defer close(r.done)

	ch, err := r.discovery.Watch(ctx, r.target)
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), slog.String("name", r.target))
		return
	}
	for {
		select {
		case <-ctx.Done():
			return
		case endpoint, ok := <-ch:
			if !ok {
				return
			}
			addr := endpoint.Address()
			if addr == "" {
				continue
			}
			switch endpoint.State {
			case StateUp:
				r.ups.Store(addr, struct{}{})
			case StateDown:
				r.ups.Delete(addr)
			default:
				slog.ErrorContext(ctx, "unknown endpoint state", slog.Int("state", int(endpoint.State)), slog.String("name", r.target), slog.String("host", endpoint.Host))
			}
			if err = r.update(); err != nil {
				slog.ErrorContext(ctx, "update conn state failed", slog.Any("error", err), slog.String("name", r.target), slog.String("host", endpoint.Host))
			}
		}
	}
}

func (r *resolver) update() error {
	var addrs []gresolver.Address
	r.ups.Range(func(key, value any) bool {
		k, _ := key.(string)
		addrs = append(addrs, gresolver.Address{Addr: k})
		return true
	})
	return r.cc.UpdateState(gresolver.State{Addresses: addrs})
}

// ResolveNow is a no-op here.
// It's just a hint, resolver can ignore this if it's not necessary.
func (r *resolver) ResolveNow(gresolver.ResolveNowOptions) {}

func (r *resolver) Close() {
	r.cancel()
	<-r.done
}
