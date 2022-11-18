package dlock

import (
	"context"
	"net/http"
	"time"

	"github.com/bufbuild/connect-go"
	lockv1 "github.com/kerraform/kerranamodb/internal/gen/lock/v1"
	"github.com/kerraform/kerranamodb/internal/gen/lock/v1/lockv1connect"
	"github.com/minio/dsync/v3"
	"go.uber.org/zap"
)

const (
	timeout = 1 * time.Second
)

type DLocker struct {
	c        lockv1connect.LockServiceClient
	endpoint string
	logger   *zap.Logger
}

type DLockerConfig struct {
	endpoint string
	logger   *zap.Logger
}

func NewDLocker(ctx context.Context, cfg *DLockerConfig) dsync.NetLocker {
	c := lockv1connect.NewLockServiceClient(
		http.DefaultClient,
		cfg.endpoint,
	)

	return &DLocker{
		c:        c,
		endpoint: cfg.endpoint,
		logger:   cfg.logger,
	}
}

func (l *DLocker) RLock(args dsync.LockArgs) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	resp, err := l.c.RLock(ctx,
		connect.NewRequest(
			&lockv1.RLockRequest{
				Uid:   args.UID,
				Table: DLockID(args.Source).Table(),
				Key:   DLockID(args.Source).Key(),
			},
		),
	)
	if err != nil {
		return false, err
	}

	return resp.Msg.GetAvailable(), nil
}

func (l *DLocker) Lock(args dsync.LockArgs) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	resp, err := l.c.Lock(ctx,
		connect.NewRequest(
			&lockv1.LockRequest{
				Uid:   args.UID,
				Table: DLockID(args.Source).Table(),
				Key:   DLockID(args.Source).Key(),
			},
		),
	)
	if err != nil {
		l.logger.Warn("failed to get availability", zap.Error(err))
		return false, err
	}

	l.logger.Info("get lock availability", zap.Bool("avilable", resp.Msg.GetAvailable()))
	return resp.Msg.GetAvailable(), nil
}

func (l *DLocker) RUnlock(args dsync.LockArgs) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	resp, err := l.c.RUnlock(ctx,
		connect.NewRequest(
			&lockv1.RUnlockRequest{
				Uid:   args.UID,
				Table: DLockID(args.Source).Table(),
				Key:   DLockID(args.Source).Key(),
			},
		),
	)
	if err != nil {
		return false, err
	}

	l.logger.Info("get lock availability", zap.Bool("avilable", resp.Msg.GetAvailable()))
	return resp.Msg.GetAvailable(), nil
}

func (l *DLocker) Unlock(args dsync.LockArgs) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	resp, err := l.c.Unlock(ctx,
		connect.NewRequest(
			&lockv1.UnlockRequest{
				Uid:   args.UID,
				Table: DLockID(args.Source).Table(),
				Key:   DLockID(args.Source).Key(),
			},
		),
	)
	if err != nil {
		return false, err
	}

	return resp.Msg.GetAvailable(), nil
}

func (l *DLocker) String() string {
	return l.endpoint
}

func (l *DLocker) Close() error {
	return nil
}
