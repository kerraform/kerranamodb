package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kerraform/kerranamodb/internal/config"
	"github.com/kerraform/kerranamodb/internal/dlock"
	"github.com/kerraform/kerranamodb/internal/driver"
	"github.com/kerraform/kerranamodb/internal/driver/local"
	"github.com/kerraform/kerranamodb/internal/driver/s3"
	"github.com/kerraform/kerranamodb/internal/http"
	server "github.com/kerraform/kerranamodb/internal/http"
	"github.com/kerraform/kerranamodb/internal/logging"
	"github.com/kerraform/kerranamodb/internal/metric"
	"github.com/kerraform/kerranamodb/internal/trace"
	v1 "github.com/kerraform/kerranamodb/internal/v1"
	"github.com/kerraform/kerranamodb/internal/version"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	otracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const (
	exitOk = iota
	exitError
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	os.Exit(exitOk)
}

func run(args []string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.Load(ctx)
	if err != nil {
		return err
	}

	logger, err := logging.NewLogger(os.Stdout, logging.Level(cfg.Log.Level), logging.Format(cfg.Log.Format))
	if err != nil {
		return err
	}

	logger = logger.With(
		zap.String("version", version.Version),
		zap.String("revision", version.Commit),
	)

	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceVersionKey.String(version.Version),
			semconv.ServiceNameKey.String(cfg.Name),
		),
	)
	if err != nil {
		logger.Error("failed to setup the otel resource", zap.Error(err))
		return err
	}

	var tp *otracesdk.TracerProvider
	if cfg.Trace.Enable {
		var sexp otracesdk.SpanExporter
		switch trace.ExporterType(cfg.Trace.Type) {
		case trace.ExporterTypeConsole:
			sexp, err = trace.NewConsoleExporter(os.Stdout)
		case trace.ExporterTypeJaeger:
			sexp, err = trace.NewJaegerExporter(cfg.Trace.Jaeger.Endpoint)
		default:
			return fmt.Errorf("trace type %s not supported", cfg.Trace.Type)
		}
		if err != nil {
			logger.Error("failed to setup the trace provider", zap.Error(err))
			return err
		}

		logger.Info("setup otel tracer", zap.String("trace", cfg.Trace.Type))
		tp = trace.NewTracer(r, sexp)
		otel.SetTracerProvider(tp)
	} else {
		logger.Debug("tracing disabled")
		tp = trace.NewTracer(r, nil)
	}
	t := tp.Tracer(cfg.Trace.Name)

	lopts := []dlock.LockOptions{dlock.WithLogger(logger.Named("dmutex"))}
	if len(cfg.Lock.Nodes) > 0 {
		lopts = append(lopts, dlock.WithStaticEndpoints(cfg.Lock.GetNodes()))
	}

	if v := cfg.Lock.ServiceDiscoveryEndpoint; v != "" {
		lopts = append(lopts, dlock.WithServiceDiscovery(v, cfg.Lock.ServiceDiscoveryNodeCount, cfg.Lock.HostIP, cfg.Lock.ServiceDiscoveryPort))
	}

	if v := cfg.Lock.ServiceDiscoveryTimeout; v != 0 {
		lopts = append(lopts, dlock.WithTimeout(time.Duration(v)*time.Second))
	}

	logger.Info("setup dlock",
		zap.Any("nodes", cfg.Lock.Nodes),
		zap.String("hostIP", cfg.Lock.HostIP),
		zap.Int("serviceDiscoveryPort", cfg.Lock.ServiceDiscoveryPort),
		zap.Int("serviceDiscoveryTimeout", cfg.Lock.ServiceDiscoveryTimeout),
		zap.String("serviceDiscoveryEndpoint", cfg.Lock.ServiceDiscoveryEndpoint),
		zap.Int("serviceDiscoveryNodeCound", cfg.Lock.ServiceDiscoveryNodeCount),
	)
	dmu, err := dlock.NewDMutex(ctx, lopts...)
	if err != nil {
		logger.Error("failed to create new lock", zap.Error(err))
		return err
	}

	logger.Info("setup backend", zap.String("backend", cfg.Backend.Type), zap.String("rootPath", cfg.Backend.RootPath))
	var d driver.Driver
	switch driver.DriverType(cfg.Backend.Type) {
	case driver.DriverTypeS3:
		d, err = s3.NewDriver(logger, &s3.DriverOpts{
			AccessKey:    cfg.Backend.S3.AccessKey,
			Bucket:       cfg.Backend.S3.Bucket,
			Endpoint:     cfg.Backend.S3.Endpoint,
			SecretKey:    cfg.Backend.S3.SecretKey,
			Tracer:       t,
			UsePathStyle: cfg.Backend.S3.UsePathStyle,
		})

		if err != nil {
			return err
		}
	case driver.DriverTypeLocal:
		d = local.NewDriver(&local.DriverConfig{
			Logger:   logger,
			Tracer:   t,
			RootPath: cfg.Backend.RootPath,
		})
	default:
		return fmt.Errorf("backend type %s not supported", cfg.Backend.Type)
	}

	metrics := metric.New(logger, d)

	wg, ctx := errgroup.WithContext(ctx)
	v1 := v1.New(&v1.HandlerConfig{
		Dmu:    dmu,
		Driver: d,
		Logger: logger,
	})

	httpSvr := http.NewServer(&server.ServerConfig{
		Dmu:    dmu,
		Driver: d,
		Logger: logger,
		Metric: metrics,
		Tracer: t,
		V1:     v1,
	})

	httpConn, err := net.Listen("tcp", cfg.HTTPAddress())
	if err != nil {
		return err
	}

	logger.Info("http server started", zap.Int("port", cfg.HTTPPort))
	wg.Go(func() error {
		return httpSvr.Serve(ctx, httpConn)
	})

	grpcSvc := dlock.NewLockService(&dlock.LockServiceOptions{
		Port:   cfg.GRPCPort,
		Logger: logger,
		Dmu:    dmu,
	})
	logger.Info("grpc server started", zap.Int("port", cfg.GRPCPort))
	wg.Go(func() error {
		return grpcSvc.Serve()
	})

	wg.Go(func() error {
		return dmu.Connect(ctx)
	})

	wg.Go(func() error {
		return errors.New("hoho")
	})

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, os.Interrupt)
	select {
	case v := <-sigCh:
		logger.Info("received signal %d", zap.String("signal", v.String()))
	case <-ctx.Done():
		return ctx.Err()
	}

	// Context for shutdown
	newCtx := context.Background()
	if err := httpSvr.Shutdown(newCtx); err != nil {
		logger.Error("failed to graceful shutdown server", zap.Error(err))
		return err
	}

	if tp != nil {
		if err := tp.Shutdown(newCtx); err != nil {
			logger.Error("failed to shutdown trace provider", zap.Error(err))
			return err
		}
	}

	return nil
}
