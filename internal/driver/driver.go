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
	DeleteLock(context.Context, string, id.LockID) error
	HasLock(context.Context, string, id.LockID) (bool, error)
	GetLock(context.Context, string, id.LockID) (Info, error)
	SaveLock(context.Context, string, id.LockID, Info) error
}

type Info string
