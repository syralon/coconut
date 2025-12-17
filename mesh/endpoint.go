package mesh

import (
	"context"
	"net"
	"strconv"
)

type State int

const (
	StateUnknown State = iota
	StateUp
	StateDown
)

type Scheme string

const (
	GRPC  = Scheme("grpc")
	HTTP  = Scheme("http")
	HTTPS = Scheme("https")
)

type Endpoint struct {
	Name     string
	Host     string
	Port     uint16
	Metadata Metadata
	State    State
	Scheme   Scheme
	ID       string
}

func (e Endpoint) Address() string {
	return net.JoinHostPort(e.Host, strconv.Itoa(int(e.Port)))
}

type Endpoints []*Endpoint

func (e Endpoints) Balance(ctx context.Context, balance Balancer) Endpoints {
	endpoints := make(Endpoints, 0, len(e))
	for _, endpoint := range e {
		if balance(ctx, endpoint) {
			endpoints = append(endpoints, endpoint)
		}
	}
	return endpoints
}

type Receipt struct{}
