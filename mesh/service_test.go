package mesh

import (
	"context"
	"slices"
	"sync"
	"testing"
	"time"
)

type inMemory struct {
	lock sync.RWMutex
	m    map[string]Endpoints

	watchers map[string]chan *Endpoint
}

func newInMemory() *inMemory {
	return &inMemory{
		m:        map[string]Endpoints{},
		watchers: map[string]chan *Endpoint{},
	}
}

func (m *inMemory) Register(ctx context.Context, endpoint *Endpoint) (*Receipt, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	endpoint.State = StateUp
	m.m[endpoint.Name] = append(m.m[endpoint.Name], endpoint)
	if ch := m.watchers[endpoint.Name]; ch != nil {
		ch <- endpoint
	}

	go m.heartbeat(ctx, endpoint)

	return &Receipt{}, nil
}

func (m *inMemory) Deregister(ctx context.Context, endpoint *Endpoint) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	endpoint.State = StateDown
	slices.DeleteFunc(m.m[endpoint.Name], func(e *Endpoint) bool { return e.Host == endpoint.Host })
	if ch := m.watchers[endpoint.Name]; ch != nil {
		ch <- endpoint
	}
	return nil
}

func (m *inMemory) heartbeat(ctx context.Context, endpoint *Endpoint) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(5 * time.Second):
			_, _ = m.Register(ctx, endpoint)
		}
	}
}

func (m *inMemory) Discover(ctx context.Context, name string) (Endpoints, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	endpoints := m.m[name]
	return endpoints, nil
}

func (m *inMemory) Watch(_ context.Context, name string) (<-chan *Endpoint, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	ch := make(chan *Endpoint)
	m.watchers[name] = ch
	return ch, nil
}

func TestRegistryDiscovery(t *testing.T) {
	ctx := context.Background()
	m := newInMemory()
	endpoint := &Endpoint{
		Name: "example",
		Host: "127.0.0.1",
		Port: 8000,
	}
	if _, err := m.Register(ctx, endpoint); err != nil {
		t.Error(err)
		return
	}

	endpoints, err := m.Discover(ctx, "example")
	if err != nil {
		t.Error(err)
	}
	for _, end := range endpoints {
		t.Log(end.Name, end.Address())
	}
}
