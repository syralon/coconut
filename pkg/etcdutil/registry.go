package etcdutil

import (
	"context"
	"fmt"

	"github.com/syralon/coconut/mesh"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
)

type Registry struct {
	client *clientv3.Client
}

func NewRegistry(client *clientv3.Client) *Registry {
	return &Registry{client: client}
}

func (r *Registry) Register(ctx context.Context, endpoint *mesh.Endpoint) (*mesh.Receipt, error) {
	em, err := endpoints.NewManager(r.client, endpoint.Name)
	if err != nil {
		return nil, err
	}
	lease, err := r.client.Grant(ctx, 30)
	if err != nil {
		return nil, err
	}
	if _, err = r.client.KeepAlive(ctx, lease.ID); err != nil {
		return nil, err
	}
	id := fmt.Sprintf("%s_%s", endpoint.Scheme, endpoint.Name)
	if err = em.AddEndpoint(ctx, id, endpoints.Endpoint{
		Addr:     endpoint.Address(),
		Metadata: endpoint.Metadata,
	}, clientv3.WithLease(lease.ID)); err != nil {
		return nil, err
	}
	return &mesh.Receipt{}, nil
}

func (r *Registry) Deregister(ctx context.Context, endpoint *mesh.Endpoint) error {
	em, err := endpoints.NewManager(r.client, endpoint.Name)
	if err != nil {
		return err
	}
	id := fmt.Sprintf("%s_%s", endpoint.Scheme, endpoint.Name)
	return em.DeleteEndpoint(ctx, id)
}
