package middleware

import (
	"net/http"

	"go.opentelemetry.io/otel/trace"
)

const (
	spanName = "server"
)

func NewTrace(t trace.Tracer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, span := t.Start(r.Context(), spanName)
			defer span.End()
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
