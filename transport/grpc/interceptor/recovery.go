package interceptor

import (
	"context"
	stderrors "errors"

	"google.golang.org/grpc"

	"github.com/syralon/coconut/proto/syralon/coconut/errors"
)

func RecoveryUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		defer func() { err = stderrors.Join(err, errors.Recovery(recover())) }()
		return handler(ctx, req)
	}
}

func RecoveryStreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() { err = stderrors.Join(err, errors.Recovery(recover())) }()
		return handler(srv, ss)
	}
}

func RecoveryUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		defer func() { err = stderrors.Join(err, errors.Recovery(recover())) }()
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func RecoveryStreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (
		stream grpc.ClientStream, err error,
	) {
		defer func() { err = stderrors.Join(err, errors.Recovery(recover())) }()
		return streamer(ctx, desc, cc, method, opts...)
	}
}
