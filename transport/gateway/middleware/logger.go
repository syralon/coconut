package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/syralon/coconut/toolkit/xio"
)

const maxRecordSize = 5 * 1024 * 1024 // 5 mb

type responseRecorder struct {
	http.ResponseWriter

	begin  time.Time
	path   string
	method string
	status int
	size   int
	buf    *xio.LimitBuffer
}

func (r *responseRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *responseRecorder) Write(data []byte) (int, error) {
	if r.buf != nil {
		_, _ = r.buf.Write(data)
	}
	size, err := r.ResponseWriter.Write(data)
	r.size += size
	return size, err
}

func (r *responseRecorder) Status() int {
	return r.status
}

func (r *responseRecorder) Size() int {
	return r.size
}

func (r *responseRecorder) Attrs() []any {
	attrs := []any{
		slog.String("method", r.method),
		slog.String("path", r.path),
		slog.Duration("duration", time.Since(r.begin)),
		slog.Int("status", r.status),
		slog.Int("size", r.size),
	}
	if r.status != http.StatusOK {
		attrs = append(attrs, slog.String("response", string(r.buf.Bytes())))
	}
	return attrs
}

type ResponseRecordWriter interface {
	http.ResponseWriter
	Status() int
	Size() int
	Attrs() []any
}

func NewRecorder(w http.ResponseWriter, r *http.Request) ResponseRecordWriter {
	re := &responseRecorder{
		ResponseWriter: w,
		begin:          time.Now(),
		path:           r.URL.Path,
		method:         r.Method,
		status:         http.StatusOK,
		size:           0,
		buf:            xio.NewLimitBuffer(maxRecordSize),
	}
	return re
}

func Logger() runtime.Middleware {
	return func(next runtime.HandlerFunc) runtime.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			rewrite := NewRecorder(w, r)
			defer func() {
				slog.InfoContext(r.Context(), r.URL.Path, rewrite.Attrs()...)
			}()
			next(rewrite, r, pathParams)
		}
	}
}
