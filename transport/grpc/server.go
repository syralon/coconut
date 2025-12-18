package grpc

import (
	"context"
	"net"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"

	"github.com/syralon/coconut/mesh"
)

type RegisterFunc func(srv *grpc.Server)

type Server struct {
	c *Config

	srv *grpc.Server

	unaryInterceptors  []grpc.UnaryServerInterceptor
	streamInterceptors []grpc.StreamServerInterceptor
	options            []grpc.ServerOption

	registers []RegisterFunc

	endpoint *mesh.Endpoint
}

func NewServer(c *Config) *Server {
	s := &Server{
		c:        c,
		endpoint: c.Endpoint(),
	}
	s.endpoint.Scheme = mesh.GRPC
	return s
}

func (s *Server) Name() string {
	return "grpc"
}

func (s *Server) Serve(_ context.Context) error {
	options := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(s.unaryInterceptors...),
		grpc.ChainStreamInterceptor(s.streamInterceptors...),
	}
	if s.c.Timeout > 0 {
		options = append(options, grpc.ConnectionTimeout(time.Duration(s.c.Timeout)))
	}
	options = append(options, s.options...)
	s.srv = grpc.NewServer(options...)
	for _, register := range s.registers {
		register(s.srv)
	}
	listener, err := net.Listen("tcp", s.c.ListenOn)
	if err != nil {
		return err
	}
	return s.srv.Serve(listener)
}

func (s *Server) Shutdown(_ context.Context) error {
	if s.srv != nil {
		s.srv.GracefulStop()
	}
	return nil
}

func (s *Server) Endpoint() (*mesh.Endpoint, bool) {
	return s.endpoint, true
}

func (s *Server) Register(registers ...RegisterFunc) *Server {
	s.registers = append(s.registers, registers...)
	return s
}

func (s *Server) WithUnaryInterceptor(interceptors ...grpc.UnaryServerInterceptor) *Server {
	s.unaryInterceptors = append(s.unaryInterceptors, interceptors...)
	return s
}

func (s *Server) WithStreamInterceptor(interceptors ...grpc.StreamServerInterceptor) *Server {
	s.streamInterceptors = append(s.streamInterceptors, interceptors...)
	return s
}

func (s *Server) WithMetadata(md mesh.Metadata) *Server {
	s.endpoint.Metadata = md
	return s
}

func (s *Server) WithOTELHandler(options ...otelgrpc.Option) {
	s.options = append(s.options, grpc.StatsHandler(otelgrpc.NewServerHandler(options...)))
}

func (s *Server) WithOptions(options ...grpc.ServerOption) *Server {
	s.options = append(s.options, options...)
	return s
}
