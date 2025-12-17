package etcdutil

import (
	"context"
	"net"
	"strconv"
	"strings"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"

	"github.com/syralon/coconut/mesh"
)

type Discovery struct {
	client *clientv3.Client
}

func (d *Discovery) Discover(ctx context.Context, name string) (mesh.Endpoints, error) {
	em, err := endpoints.NewManager(d.client, name)
	if err != nil {
		return nil, err
	}
	list, err := em.List(ctx)
	if err != nil {
		return nil, err
	}
	eds := make(mesh.Endpoints, 0, len(list))
	for id, endpoint := range list {
		eds = append(eds, d.build(name, id, endpoint, mesh.StateUp))
	}
	return eds, nil
}

func (d *Discovery) Watch(ctx context.Context, name string) (<-chan *mesh.Endpoint, error) {
	em, err := endpoints.NewManager(d.client, name)
	if err != nil {
		return nil, err
	}
	wch, err := em.NewWatchChannel(ctx)
	if err != nil {
		return nil, err
	}
	ch := make(chan *mesh.Endpoint)
	go d.watch(ctx, name, wch, ch)
	return ch, nil
}

func (d *Discovery) watch(ctx context.Context, name string, wch endpoints.WatchChannel, ch chan<- *mesh.Endpoint) {
	defer close(ch)
	for {
		select {
		case <-ctx.Done():
			return
		case updates, ok := <-wch:
			if !ok {
				return
			}
			for _, update := range updates {
				ed := d.build(name, update.Key, update.Endpoint, mesh.StateDown)
				switch update.Op {
				case endpoints.Add:
					ed.State = mesh.StateUp
				case endpoints.Delete:
					ed.State = mesh.StateDown
				}
				ch <- ed
			}
		}
	}
}

func (d *Discovery) build(name string, id string, endpoint endpoints.Endpoint, state mesh.State) *mesh.Endpoint {
	var md map[string][]string
	if endpoint.Metadata != nil {
		if val, ok := endpoint.Metadata.(map[string][]string); ok {
			md = val
		}
	}
	var scheme mesh.Scheme
	if n := strings.Index(id, "_"); n > 0 {
		scheme = mesh.Scheme(id[:n])
		id = id[n+1:]
	}
	host, port, _ := net.SplitHostPort(endpoint.Addr)
	portInt, _ := strconv.Atoi(port)
	ed := &mesh.Endpoint{
		Name:     name,
		Host:     host,
		Port:     uint16(portInt),
		Metadata: md,
		State:    state,
		Scheme:   scheme,
		ID:       id,
	}
	return ed
}
