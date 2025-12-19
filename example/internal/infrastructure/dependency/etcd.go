package dependency

import (
	"github.com/syralon/coconut/example/internal/config"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func NewETCDClient(c *config.Config) (*clientv3.Client, func(), error) {
	client, err := c.ETCD.NewClient()
	if err != nil {
		return nil, nil, err
	}
	return client, func() { _ = client.Close() }, err
}
