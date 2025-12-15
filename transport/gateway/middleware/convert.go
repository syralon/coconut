package middleware

import (
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func convert(std func(http.Handler) http.Handler) runtime.Middleware {
	return func(next runtime.HandlerFunc) runtime.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			h := std(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				next(writer, request, pathParams)
			}))
			h.ServeHTTP(w, r)
		}
	}
}
