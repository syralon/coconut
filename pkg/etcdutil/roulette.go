package etcdutil

import (
	"context"
	"fmt"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type RouletteOption func(*Roulette)

func WithMaxID(maxID int) RouletteOption {
	return func(m *Roulette) {
		m.maxID = maxID
	}
}

type Roulette struct {
	name    string
	client  *clientv3.Client
	leaseID clientv3.LeaseID
	maxID   int
}

func NewRoulette(name string, client *clientv3.Client, options ...RouletteOption) *Roulette {
	const defaultMaxID = 4096
	r := &Roulette{
		name:   name,
		client: client,
		maxID:  defaultMaxID,
	}
	for _, option := range options {
		option(r)
	}
	return r
}

func (r *Roulette) Allocate(ctx context.Context) (int, error) {
	lease, err := r.client.Grant(ctx, 30) // 30s TTL
	if err != nil {
		return 0, err
	}
	r.leaseID = lease.ID
	_, err = r.client.KeepAlive(ctx, lease.ID)
	if err != nil {
		return 0, err
	}
	for i := 1; i <= r.maxID; i++ {
		key := fmt.Sprintf("/%s/workers/%d", r.name, i)
		txn := r.client.Txn(ctx).
			If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)). //
			Then(clientv3.OpPut(key, "1", clientv3.WithLease(lease.ID)))
		resp, err := txn.Commit()
		if err != nil {
			return 0, err
		}
		if resp.Succeeded {
			return i, nil
		}
	}
	return 0, fmt.Errorf("no available workerId (0-%d) for service %s", r.maxID, r.name)
}

func (r *Roulette) Release(ctx context.Context) error {
	_, err := r.client.Revoke(ctx, r.leaseID)
	return err
}
