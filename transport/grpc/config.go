package grpc

import (
	"github.com/syralon/coconut/transport/common"
	"google.golang.org/grpc"
)

type Config struct {
	common.Config
}

type ClientConfig struct {
	Target string `json:"target"  yaml:"target"`
	Secure bool   `json:"secure"  yaml:"secure"`
}

type ServiceClientConfig[T any] struct {
	ClientConfig
}

// ConnectFunc create a new service client from current client config.
// The 'T' is a grpc client interface type.
// The 'fn' is the grpc generate 'NewXXXClient' function.
// e.g.:
//
//	var service examplepb.ExampleServiceClient
//	var config = ServiceClientConfig[examplepb.ExampleServiceClient]
//	var cc = NewConnector()
//	cc.Connects(
//	  config.Connect(&service, examplepb.NewExampleServiceClient),
//	  // ... others service client
//	)
func (s *ServiceClientConfig[T]) ConnectFunc(service *T, fn func(grpc.ClientConnInterface) T) func(*Connector) error {
	return ConnectFunc[T](service, &s.ClientConfig, fn)
}

func (s *ServiceClientConfig[T]) Connect(connector *Connector, service *T, fn func(grpc.ClientConnInterface) T) error {
	return Connect(connector, service, &s.ClientConfig, fn)
}
