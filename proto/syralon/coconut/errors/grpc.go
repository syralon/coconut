package errors

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func error2status(err error) error {
	if _, ok := status.FromError(err); ok || err == nil {
		return err
	}
	if ee, ok := FromError(err); ok {
		err = ee.Status().Err()
	}
	return err
}

func status2error(err error) error {
	if err == nil {
		return nil
	}
	st, ok := status.FromError(err)
	if !ok {
		return err
	}
	if ee, ok := FromStatus(st); ok {
		return ee
	}
	return err
}

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		defer func() { err = error2status(err) }()
		return handler(ctx, req)
	}
}

type errorStreamServer struct {
	grpc.ServerStream
}

func (s *errorStreamServer) SendMsg(m any) (err error) {
	defer func() { err = error2status(err) }()
	return s.ServerStream.SendMsg(m)
}

func (s *errorStreamServer) RecvMsg(m any) (err error) {
	defer func() { err = error2status(err) }()
	return s.ServerStream.RecvMsg(m)
}

func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		return handler(srv, &errorStreamServer{ServerStream: ss})
	}
}

func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		defer func() { err = status2error(err) }()
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

type errorClientStream struct {
	grpc.ClientStream
}

func (s *errorClientStream) SendMsg(m any) (err error) {
	defer func() { err = status2error(err) }()
	return s.ClientStream.SendMsg(m)
}

func (s *errorClientStream) RecvMsg(m any) (err error) {
	defer func() { err = status2error(err) }()
	return s.ClientStream.RecvMsg(m)
}

func StreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		cs, err := streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			return nil, status2error(err)
		}
		return &errorClientStream{cs}, nil
	}
}
