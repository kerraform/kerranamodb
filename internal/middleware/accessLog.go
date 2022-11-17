package middleware

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kerraform/kerranamodb/internal/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type rwWrapper struct {
	rw     http.ResponseWriter
	mirror *http.Response
	closed bool
}

// newRwWrapper wraps the HTTP responseWriter for audit logging
func newRwWrapper(rw http.ResponseWriter, mirror *http.Response) *rwWrapper {
	return &rwWrapper{
		rw:     rw,
		mirror: mirror,
	}
}

func (r *rwWrapper) Header() http.Header {
	return r.rw.Header()
}

func (r *rwWrapper) Write(i []byte) (int, error) {
	r.mirror.Body = ioutil.NopCloser(bytes.NewReader(i))
	return r.rw.Write(i)
}

func (r *rwWrapper) WriteHeader(statusCode int) {
	if r.closed {
		return
	}
	r.closed = true
	r.rw.WriteHeader(statusCode)
	r.mirror.StatusCode = statusCode
}

func AccessLog(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			res := &http.Response{}
			rww := newRwWrapper(w, res)

			l := logger.With(
				flattenFields(r)...,
			)

			req := r.WithContext(context.WithValue(r.Context(), logging.Key, l))

			defer func() {
				l.Named("accessLog").Info("access to server",
					zap.Int("statusCode", res.StatusCode),
				)
			}()
			next.ServeHTTP(rww, req)
		})
	}
}

func flattenFields(r *http.Request) []zapcore.Field {
	fs := []zapcore.Field{
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.String("userAgent", r.UserAgent()),
		zap.String("contentLength", strconv.FormatInt(r.ContentLength, 10)),
		zap.String("query", r.URL.Query().Encode()),
	}
	for k, v := range mux.Vars(r) {
		fs = append(fs, zap.String(k, v))
	}

	return fs
}
