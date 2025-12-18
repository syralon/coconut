package etcdutil

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"

	clientv3 "go.etcd.io/etcd/client/v3"
	"sigs.k8s.io/yaml"

	"github.com/syralon/coconut/configuration"
)

type etcdConfigDriver struct{}

func NewConfigDriver() configuration.Driver {
	return &etcdConfigDriver{}
}

func (c *etcdConfigDriver) Build(_ context.Context, script string) (configuration.Reader, error) {
	var config *clientv3.Config
	var err error
	if strings.HasPrefix(script, "etcd://") {
		config, err = c.configFromURL(script)
	} else {
		config, err = c.configFromLocalFile(script)
	}
	if err != nil {
		return nil, err
	}
	client, err := clientv3.New(*config)
	if err != nil {
		return nil, err
	}
	return &etcdConfigReader{client: client}, nil
}

func (c *etcdConfigDriver) configFromURL(u string) (*clientv3.Config, error) {
	uu, err := url.Parse(u)
	if err != nil {
		return nil, err
	}
	config := &clientv3.Config{Endpoints: strings.Split(uu.Host, ",")}
	if uu.User != nil {
		config.Username = uu.User.Username()
		config.Password, _ = uu.User.Password()
	}
	return config, nil
}

func (c *etcdConfigDriver) configFromLocalFile(filename string) (*clientv3.Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	switch strings.ToLower(path.Ext(filename)) {
	case ".yaml", ".yml":
		err = yaml.Unmarshal(data, cfg)
	case ".json":
		err = json.Unmarshal(data, cfg)
	}
	if err != nil {
		return nil, err
	}
	cc := cfg.Config()
	return &cc, nil
}

type etcdConfigReader struct {
	client *clientv3.Client
}

func (c *etcdConfigReader) Read(ctx context.Context, key string) ([]byte, error) {
	data, err := c.client.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	for _, kv := range data.Kvs {
		if string(kv.Key) == key {
			return kv.Value, nil
		}
	}
	return nil, fmt.Errorf("%s not found", key)
}

func (c *etcdConfigReader) Watch(ctx context.Context, key string) (<-chan *configuration.Content, error) {
	ch := make(chan *configuration.Content)
	go c.watch(ctx, c.client.Watch(ctx, key), ch)
	return ch, nil
}

func (c *etcdConfigReader) watch(ctx context.Context, wc clientv3.WatchChan, ch chan *configuration.Content) {
	defer close(ch)
	for {
		select {
		case <-ctx.Done():
			return
		case resp, ok := <-wc:
			if !ok {
				return
			}
			for _, event := range resp.Events {
				if event.Type == clientv3.EventTypeDelete {
					continue
				}
				ch <- &configuration.Content{
					Key:  string(event.Kv.Key),
					Data: event.Kv.Value,
				}
			}
		}
	}
}
