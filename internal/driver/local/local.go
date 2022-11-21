package local

import (
	"context"

	"github.com/kerraform/kerranamodb/internal/driver"
	"github.com/kerraform/kerranamodb/internal/id"
	modelv1 "github.com/kerraform/kerranamodb/internal/model/v1"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type d struct {
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
	return &d{
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

func (d *d) CreateTenant(ctx context.Context, table string, token string) error {
	return nil
}

func (d *d) GetTenant(ctx context.Context, table string) (*modelv1.Tenant, error) {
	return nil, nil
}
