package mesh

import (
	"context"
	"net"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	gresolver "google.golang.org/grpc/resolver"
	"google.golang.org/grpc/serviceconfig"
)

type conn struct {
	ch chan gresolver.State
}

var _ gresolver.ClientConn = (*conn)(nil)

func (c *conn) UpdateState(state gresolver.State) error {
	c.ch <- state
	return nil
}

func (c *conn) ReportError(error) {}

func (c *conn) NewAddress(addresses []gresolver.Address) {}

func (c *conn) ParseServiceConfig(serviceConfigJSON string) *serviceconfig.ParseResult {
	return nil
}

func TestResolver(t *testing.T) {
	memory := newInMemory()

	b := NewBuilder("grpc", memory)
	u, _ := url.Parse("grpc:///example")
	target := gresolver.Target{URL: *u}

	ctx := context.Background()
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	endpoint := &Endpoint{
		Name: "example",
		Host: "127.0.0.1",
		Port: 8000,
	}

	go func() {
		time.Sleep(1 * time.Second)
		_, err := memory.Register(ctx, endpoint)
		if err != nil {
			t.Error(err)
		}
	}()

	go func() {
		time.Sleep(2 * time.Second)
		memory.Deregister(ctx, endpoint)
	}()

	ch := make(chan gresolver.State, 3)
	cc := &conn{ch: ch}
	_, err := b.Build(target, cc, gresolver.BuildOptions{})
	if err != nil {
		t.Error(err)
	}

	expected := []gresolver.State{
		{}, // before register, got empty addresses
		{Addresses: []gresolver.Address{{Addr: endpoint.Address()}}}, // after registered, got one new address
		{}, // after deregistered, got empty addresses
	}
	for i := 0; i < 3; i++ {
		select {
		case state := <-ch:
			t.Log(state)
			assert.Equal(t, state, expected[i])
		case <-ctx.Done():
			t.Error("timeout")
		}
	}
}

func TestServiceResolver(t *testing.T) {
}

type ResolverSuite struct {
	suite.Suite

	memory *inMemory
}

func (s *ResolverSuite) SetupSuite() {
	s.memory = newInMemory()

	srv := grpc.NewServer()
	grpc_health_v1.RegisterHealthServer(srv, health.NewServer())

	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		s.T().Fatal(err)
	}
	ech := make(chan error, 1)
	go func() { ech <- srv.Serve(listener) }()
	select {
	case err = <-ech:
		s.T().Fatal(err)
	case <-time.After(100 * time.Millisecond):
	}
	if _, err = s.memory.Register(context.Background(), &Endpoint{
		Name: "example",
		Host: "127.0.0.1",
		Port: 8000,
	}); err != nil {
		s.T().Fatal(err)
	}
}

func (s *ResolverSuite) TestResolve() {
	cc, err := grpc.NewClient(
		"example:///example",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithResolvers(NewBuilder("example", s.memory)),
	)
	if err != nil {
		s.T().Fatal(err)
	}
	client := grpc_health_v1.NewHealthClient(cc)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := client.List(ctx, &grpc_health_v1.HealthListRequest{})
	if err != nil {
		s.Error(err)
		return
	}
	assert.Equal(s.T(), 1, len(resp.GetStatuses()))
	for k, v := range resp.GetStatuses() {
		s.T().Log(k, v.String())
		assert.Equal(s.T(), v.Status, grpc_health_v1.HealthCheckResponse_SERVING)
	}
}

func TestResolverSuite(t *testing.T) {
	suite.Run(t, new(ResolverSuite))
}
