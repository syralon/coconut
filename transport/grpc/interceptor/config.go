package interceptor

import "google.golang.org/grpc"

type Config struct {
	Logger *LoggerConfig
}

func (c *Config) UnaryServerInterceptors() (interceptors []grpc.UnaryServerInterceptor) {
	if c == nil {
		return nil
	}
	if c.Logger != nil {
		interceptors = append(interceptors, LoggerUnaryServerInterceptor(c.Logger))
	}
	return interceptors
}

func (c *Config) StreamServerInterceptors() (interceptors []grpc.StreamServerInterceptor) {
	if c == nil {
		return nil
	}

	return interceptors
}

func (c *Config) UnaryClientInterceptors() (interceptors []grpc.UnaryClientInterceptor) {
	if c == nil {
		return nil
	}
	if c.Logger != nil {
		interceptors = append(interceptors, LoggerUnaryClientInterceptor(c.Logger))
	}
	return interceptors
}

func (c *Config) StreamClientInterceptors() (interceptors []grpc.StreamClientInterceptor) {
	if c == nil {
		return nil
	}

	return interceptors
}
