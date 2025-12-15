package interceptor

import "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

var (
	NewServerHandler = otelgrpc.NewServerHandler
	NewClientHandler = otelgrpc.NewClientHandler
)
