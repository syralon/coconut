package mesh

import "context"

type Registry interface {
	Register(ctx context.Context, endpoint *Endpoint) (*Receipt, error)
	Deregister(ctx context.Context, endpoint *Endpoint) error
}

type Discovery interface {
	Discover(ctx context.Context, name string) (Endpoints, error)
	Watch(ctx context.Context, name string) (<-chan *Endpoint, error)
}

type Balancer func(ctx context.Context, endpoint *Endpoint) bool
