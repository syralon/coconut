package middleware

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func OTEL() runtime.Middleware {
	return convert(otelhttp.NewMiddleware("coconut-http-server"))
}
