package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewClient(c *ClientConfig, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	if !c.Secure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	return grpc.NewClient(c.Target, opts...)
}

type Connector struct {
	clients map[string]*grpc.ClientConn
}

func NewConnector() *Connector {
	return &Connector{clients: make(map[string]*grpc.ClientConn)}
}

func (c *Connector) Close() error {
	for _, client := range c.clients {
		_ = client.Close()
	}
	return nil
}

func (c *Connector) NewClient(config *ClientConfig, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	if cc, ok := c.clients[config.Target]; ok {
		return cc, nil
	}
	if !config.Secure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	cc, err := NewClient(config, opts...)
	if err != nil {
		return nil, err
	}
	c.clients[config.Target] = cc
	return cc, nil
}

func (c *Connector) Connects(fns ...func(connector *Connector) error) error {
	for _, fn := range fns {
		if err := fn(c); err != nil {
			return err
		}
	}
	return nil
}

func ConnectFunc[T any](service *T, config *ClientConfig, fn func(grpc.ClientConnInterface) T, options ...grpc.DialOption) func(*Connector) error {
	return func(connector *Connector) error {
		return Connect(connector, service, config, fn, options...)
	}
}

func Connect[T any](connector *Connector, service *T, config *ClientConfig, fn func(grpc.ClientConnInterface) T, options ...grpc.DialOption) (err error) {
	cc, err := connector.NewClient(config, options...)
	if err != nil {
		return err
	}
	*service = fn(cc)
	return nil
}
