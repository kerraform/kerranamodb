package local

import (
	"github.com/kerraform/kerranamodb/internal/driver"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type local struct {
	rootPath string
	logger   *zap.Logger
	tracer   trace.Tracer
}

type DriverConfig struct {
	RootPath string
	Logger   *zap.Logger
	Tracer   trace.Tracer
}

func NewDriver(cfg *DriverConfig) driver.Driver {
	return &local{
		logger:   cfg.Logger,
		rootPath: cfg.RootPath,
		tracer:   cfg.Tracer,
	}
}
