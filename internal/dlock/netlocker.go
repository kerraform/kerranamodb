package dlock

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/kerraform/kerranamodb/internal/gen/lock/v1/lockv1connect"
	"github.com/minio/dsync/v3"
)

type DLocker struct {
	c lockv1connect.LockServiceClient
}

func NewDLocker(ctx context.Context, endpoint string) dsync.NetLocker {
	c := lockv1connect.NewLockServiceClient(
		http.DefaultClient,
		endpoint,
	)

	return &DLocker{
		c: c,
	}
}

func (l *DLocker) RLock(args dsync.LockArgs) (bool, error) {
	fmt.Println("rlock 1")
	return true, nil
}

func (l *DLocker) Lock(args dsync.LockArgs) (bool, error) {
	fmt.Println("rlock 2")
	time.Sleep(30 * time.Second)
	return true, nil
}

func (l *DLocker) RUnlock(args dsync.LockArgs) (bool, error) {
	fmt.Println("rlock 3")
	return true, nil
}

func (l *DLocker) Unlock(args dsync.LockArgs) (bool, error) {
	fmt.Println("rlock 4")
	return true, nil
}

func (l *DLocker) String() string {
	return ""
}

func (l *DLocker) Close() error {
	return nil
}
