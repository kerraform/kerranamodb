package driver

import "errors"

var (
	ErrLockNotFound   = errors.New("lock not found")
	ErrTenantNotFound = errors.New("tenant exists")
)
