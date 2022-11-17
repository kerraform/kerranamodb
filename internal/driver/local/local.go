package local

import (
	"context"

	"github.com/kerraform/kerranamodb/internal/dlock"
	"github.com/kerraform/kerranamodb/internal/driver"
	"github.com/kerraform/kerranamodb/internal/id"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type d struct {
	dmu      *dlock.DMutex
	rootPath string
	logger   *zap.Logger
	tracer   trace.Tracer
}

type DriverConfig struct {
	Dmu      *dlock.DMutex
	RootPath string
	Logger   *zap.Logger
	Tracer   trace.Tracer
}

func NewDriver(cfg *DriverConfig) driver.Driver {
	return &d{
		dmu:      cfg.Dmu,
		logger:   cfg.Logger,
		rootPath: cfg.RootPath,
		tracer:   cfg.Tracer,
	}
}

func (d *d) DeleteLock(ctx context.Context, table string, lid id.LockID) error {
	return nil
}

func (d *d) HasLock(ctx context.Context, table string, lid id.LockID) (bool, error) {
	return true, nil
}

func (d *d) GetLock(ctx context.Context, table string, lid id.LockID) (driver.Info, error) {
	return "", nil
}

func (d *d) SaveLock(ctx context.Context, table string, lid id.LockID, info driver.Info) error {
	return nil
}
