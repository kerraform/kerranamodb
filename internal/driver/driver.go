package driver

import (
	"context"

	"github.com/kerraform/kerranamodb/internal/id"
	modelv1 "github.com/kerraform/kerranamodb/internal/model/v1"
)

type DriverType string

const (
	DriverTypeLocal DriverType = "local"
	DriverTypeS3    DriverType = "s3"
)

const (
	TokenFile = "token"
)

type Driver interface {
	// Lock
	DeleteLock(context.Context, string, id.LockID) error
	HasLock(context.Context, string, id.LockID) (bool, error)
	GetLock(context.Context, string, id.LockID) (Info, error)
	SaveLock(context.Context, string, id.LockID, Info) error

	// Tenant
	CreateTenant(context.Context, string, string) error
	GetTenant(context.Context, string) (*modelv1.Tenant, error)
}

type Info string
