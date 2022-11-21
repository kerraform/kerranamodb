package http

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/kerraform/kerranamodb/internal/auth"
	"github.com/kerraform/kerranamodb/internal/dlock"
	"github.com/kerraform/kerranamodb/internal/driver"
	"github.com/kerraform/kerranamodb/internal/metric"
	"github.com/kerraform/kerranamodb/internal/middleware"
	v1 "github.com/kerraform/kerranamodb/internal/v1"
	"go.opentelemetry.io/otel/trace"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Server struct {
	auth           auth.Authenticator
	dmu            *dlock.DMutex
	driver         driver.Driver
	enableModule   bool
	enableProvider bool
	logger         *zap.Logger
	metric         *metric.RegistryMetrics
	mux            *mux.Router
	server         *http.Server
	tracer         trace.Tracer
	corsOrigin     string

	v1 *v1.Handler
}

type ServerConfig struct {
	Auth           auth.Authenticator
	Dmu            *dlock.DMutex
	Driver         driver.Driver
	EnableModule   bool
	EnableProvider bool
	Logger         *zap.Logger
	Metric         *metric.RegistryMetrics
	Tracer         trace.Tracer
	CORSOrigin     string
	V1             *v1.Handler
}

func NewServer(cfg *ServerConfig) *Server {
	s := &Server{
		auth:           cfg.Auth,
		driver:         cfg.Driver,
		dmu:            cfg.Dmu,
		enableModule:   cfg.EnableModule,
		enableProvider: cfg.EnableProvider,
		logger:         cfg.Logger,
		metric:         cfg.Metric,
		tracer:         cfg.Tracer,
		mux:            mux.NewRouter(),
		v1:             cfg.V1,
		corsOrigin:     cfg.CORSOrigin,
	}

	if cfg.Tracer != nil {
		s.mux.Use(middleware.NewTrace(s.tracer))
	}

	s.registerRegistryHandler()
	s.metric.RegisterAllMetrics()
	s.registerUtilHandler()
	s.registerMetricsHandler()

	return s
}

func (s *Server) Serve(ctx context.Context, conn net.Listener) error {
	server := &http.Server{
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      s.mux,
	}

	s.metric.Resync(ctx)
	s.server = server
	if err := server.Serve(conn); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
