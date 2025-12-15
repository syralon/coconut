package gateway

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/textproto"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/syralon/coconut/proto/syralon/coconut/errors"
)

func DefaultErrorHandler() runtime.ErrorHandlerFunc {
	return func(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
		const fallback = `{"code": 13, "message": "failed to marshal error message"}`

		ee, ok := errors.FromError(err)
		if !ok {
			ee = errors.STATUS_INTERNAL_SERVER_ERROR.ToError()
			slog.ErrorContext(ctx, err.Error())
		}
		if ee.GetStatusCode() == 0 {
			ee.StatusCode = http.StatusInternalServerError
		}

		w.Header().Del("Trailer")
		w.Header().Del("Transfer-Encoding")
		if md, ok := runtime.ServerMetadataFromContext(ctx); ok {
			writeServerMetadata(w, md)
			if requestAcceptsTrailers(r) {
				writeTrailerHeader(w, md)
				w.Header().Set("Transfer-Encoding", "chunked")
			}
		}

		w.Header().Set("Content-Type", marshaler.ContentType(ee))

		if ee.GetStatusCode() == http.StatusUnauthorized {
			w.Header().Set("WWW-Authenticate", ee.GetMessage())
		}

		buf, merr := marshaler.Marshal(ee)
		if merr != nil {
			slog.ErrorContext(ctx, "marshal error", slog.Any("err", merr))
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = io.WriteString(w, fallback)
			return
		}

		w.WriteHeader(int(ee.GetStatusCode()))
		if _, err := w.Write(buf); err != nil {
			slog.ErrorContext(ctx, "write response error", slog.Any("err", err))
		}

		if md, ok := runtime.ServerMetadataFromContext(ctx); ok && requestAcceptsTrailers(r) {
			writeTrailers(w, md)
		}
	}
}

func writeServerMetadata(w http.ResponseWriter, md runtime.ServerMetadata) {
	for k, vv := range md.HeaderMD {
		key := textproto.CanonicalMIMEHeaderKey("Grpc-Metadata-" + k)
		for _, v := range vv {
			w.Header().Add(key, v)
		}
	}
}

func requestAcceptsTrailers(r *http.Request) bool {
	te := r.Header.Get("TE")
	for _, t := range strings.Split(te, ",") {
		if strings.TrimSpace(strings.ToLower(t)) == "trailers" {
			return true
		}
	}
	return false
}

func writeTrailerHeader(w http.ResponseWriter, md runtime.ServerMetadata) {
	for k := range md.TrailerMD {
		trailerKey := textproto.CanonicalMIMEHeaderKey("Grpc-Trailer-" + k)
		w.Header().Add("Trailer", trailerKey)
	}
}

func writeTrailers(w http.ResponseWriter, md runtime.ServerMetadata) {
	for k, vv := range md.TrailerMD {
		key := textproto.CanonicalMIMEHeaderKey("Grpc-Trailer-" + k)
		for _, v := range vv {
			w.Header().Add(key, v)
		}
	}
}
