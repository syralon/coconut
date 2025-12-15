package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type recorder struct {
	begin  time.Time
	path   string
	method string
	http.ResponseWriter
	status int
	size   int
}

func (r *recorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *recorder) Write(data []byte) (int, error) {
	size, err := r.ResponseWriter.Write(data)
	r.size += size
	return size, err
}

func (r *recorder) Status() int {
	return r.status
}

func (r *recorder) Size() int {
	return r.size
}

func (r *recorder) Attrs() []any {
	return []any{
		slog.String("method", r.method),
		slog.String("path", r.path),
		slog.Duration("duration", time.Since(r.begin)),
		slog.Int("status", r.status),
		slog.Int("size", r.size),
	}
}

type ResponseRecordWriter interface {
	http.ResponseWriter
	Status() int
	Size() int
	Attrs() []any
}

func NewRecorder(w http.ResponseWriter, r *http.Request) ResponseRecordWriter {
	re := &recorder{
		begin:          time.Now(),
		path:           r.URL.Path,
		method:         r.Method,
		ResponseWriter: w,
		status:         http.StatusOK,
		size:           0,
	}
	return re
}

func Logger() runtime.Middleware {
	return func(next runtime.HandlerFunc) runtime.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			rewrite := NewRecorder(w, r)
			defer func() { slog.InfoContext(r.Context(), r.URL.Path, rewrite.Attrs()...) }()
			next(rewrite, r, pathParams)
		}
	}
}
