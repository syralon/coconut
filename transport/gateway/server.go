package gateway

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"github.com/syralon/coconut/mesh"
)

type (
	EndpointRegister func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)
	ConnRegister     func(ctx context.Context, mux *runtime.ServeMux, cc *grpc.ClientConn) error
	Register         func(ctx context.Context, mux *runtime.ServeMux) error
)

type Server struct {
	c *Config

	srv *http.Server

	options     []runtime.ServeMuxOption
	dialOptions []grpc.DialOption
	registers   []Register

	endpoint *mesh.Endpoint
}

func NewServer(c *Config) *Server {
	s := &Server{
		c:   c,
		srv: &http.Server{Addr: c.ListenOn},
		options: []runtime.ServeMuxOption{
			runtime.WithErrorHandler(DefaultErrorHandler()),
		},
		endpoint: c.Endpoint(),
	}
	if s.c.TLS.KeyFile != "" && s.c.TLS.CertFile != "" {
		s.endpoint.Scheme = mesh.HTTPS
	} else {
		s.endpoint.Scheme = mesh.HTTP
	}
	return s
}

func (s *Server) WithOptions(opts ...runtime.ServeMuxOption) *Server {
	s.options = append(s.options, opts...)
	return s
}

func (s *Server) WithDialOptions(opts ...grpc.DialOption) *Server {
	s.dialOptions = append(s.dialOptions, opts...)
	return s
}

func (s *Server) WithMetadata(md mesh.Metadata) *Server {
	s.endpoint.Metadata = md
	return s
}

func (s *Server) Register(fns ...Register) *Server {
	s.registers = append(s.registers, fns...)
	return s
}

func (s *Server) RegisterEndpoint(endpoint string, fns ...EndpointRegister) *Server {
	s.registers = append(s.registers, func(ctx context.Context, mux *runtime.ServeMux) error {
		for _, fn := range fns {
			if err := fn(ctx, mux, endpoint, s.dialOptions); err != nil {
				return err
			}
		}
		return nil
	})
	return s
}

func (s *Server) RegisterConn(cc *grpc.ClientConn, fns ...ConnRegister) *Server {
	s.registers = append(s.registers, func(ctx context.Context, mux *runtime.ServeMux) error {
		for _, fn := range fns {
			if err := fn(ctx, mux, cc); err != nil {
				return err
			}
		}
		return nil
	})
	return s
}

func (s *Server) Serve(ctx context.Context) error {
	mux := runtime.NewServeMux(s.options...)
	for _, register := range s.registers {
		if err := register(ctx, mux); err != nil {
			return err
		}
	}
	s.srv.Handler = mux
	if s.c.TLS.CertFile != "" && s.c.TLS.KeyFile != "" {
		return s.srv.ListenAndServeTLS(s.c.TLS.CertFile, s.c.TLS.KeyFile)
	}
	return s.srv.ListenAndServe()
}

func (s *Server) Name() string {
	return "gateway"
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func (s *Server) Endpoint() (*mesh.Endpoint, bool) {
	return s.endpoint, true
}
