package etcdutil

import (
	"context"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func keepalive(ctx context.Context, client *clientv3.Client, leaseID clientv3.LeaseID) error {
	ch, err := client.KeepAlive(ctx, leaseID)
	if err != nil {
		return err
	}
	go func() {
		for range ch {
		}
	}()
	return nil
}
