package driver

import (
	"context"

	"github.com/kerraform/kerranamodb/internal/id"
)

type DriverType string

const (
	DriverTypeLocal DriverType = "local"
	DriverTypeS3    DriverType = "s3"
)

type Driver interface {
	// Lock
	DeleteLock(context.Context, string, id.LockID) error
	HasLock(context.Context, string, id.LockID) (bool, error)
	GetLock(context.Context, string, id.LockID) (Info, error)
	SaveLock(context.Context, string, id.LockID, Info) error

	// Tenant
	CreateTenant(context.Context, string) error
	GetTenant(context.Context, string) error
}

type Info string
