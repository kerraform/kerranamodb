package middleware

import (
	"net/http"

	"github.com/kerraform/kerranamodb/internal/metric"
)

func AccessMetric(m *metric.RegistryMetrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			res := &http.Response{}
			rww := newRwWrapper(w, res)

			defer func() {
				m.IncrementHTTPRequestTotal(res.StatusCode, r.Method, r.URL.Path)
			}()
			next.ServeHTTP(rww, r)
		})
	}
}
