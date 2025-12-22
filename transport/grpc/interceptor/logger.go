package interceptor

import (
	"context"
	"encoding/json"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type LoggerStrategy string

const (
	Always    LoggerStrategy = "always"
	OnlyError LoggerStrategy = "only_error"
)

type LoggerConfig struct {
	Enable         bool           `json:"enable"          yaml:"enable"`          //
	RecordStrategy LoggerStrategy `json:"record_strategy" yaml:"record_strategy"` // always / error_only / none
	RecordMessage  LoggerStrategy `json:"record_message"  yaml:"record_message"`  // always / error_only / none
}

func LoggerUnaryServerInterceptor(config *LoggerConfig) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (reply any, err error) {
		defer func() {
			config.log(ctx, info.FullMethod, req, reply, err)
		}()
		return handler(ctx, req)
	}
}

func LoggerUnaryClientInterceptor(config *LoggerConfig) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		defer func() {
			config.log(ctx, method, req, reply, err)
		}()
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func (c *LoggerConfig) log(ctx context.Context, method string, req, reply any, err error) {
	if !c.Enable || (err == nil && c.RecordStrategy == OnlyError) {
		return
	}
	attrs := []any{
		slog.String("method", method),
	}
	if c.RecordMessage == Always || (c.RecordMessage == OnlyError && err != nil) {
		attrs = append(attrs, slog.String("request", c.format(ctx, req)), slog.String("reply", c.format(ctx, reply)))
	}
	if err == nil {
		slog.InfoContext(ctx, method, attrs...)
	} else {
		attrs = append(attrs, slog.String("error", err.Error()))
		slog.ErrorContext(ctx, method, attrs...)
	}
}

func (c *LoggerConfig) format(ctx context.Context, m any) string {
	var data []byte
	var err error
	if msg, ok := m.(proto.Message); ok {
		data, err = protojson.Marshal(msg)
	} else {
		data, err = json.Marshal(m)
	}
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
	}
	return string(data)
}
