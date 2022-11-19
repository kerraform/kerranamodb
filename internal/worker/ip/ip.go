package ip

import (
	"context"
	"time"

	"github.com/kerraform/kerranamodb/internal/dlock"
	"github.com/kerraform/kerranamodb/internal/worker"
	"go.uber.org/zap"
)

var (
	defaultPeriod = 60
)

type syncer struct {
	dmu      *dlock.DMutex
	endpoint string
	hostIP   string
	logger   *zap.Logger
	period   int
}

var _ worker.Worker = (*syncer)(nil)

type options struct {
	logger *zap.Logger
	period int
}

type SyncerOptions func(*options)

func WithLogger(logger *zap.Logger) SyncerOptions {
	return func(o *options) {
		o.logger = logger
	}
}

func WithSyncPeriod(period int) SyncerOptions {
	return func(o *options) {
		o.period = period
	}
}

func NewSyncer(dmu *dlock.DMutex, opts ...SyncerOptions) worker.Worker {
	o := &options{
		logger: zap.NewNop(),
		period: defaultPeriod,
	}

	for _, opt := range opts {
		opt(o)
	}

	return &syncer{
		dmu: dmu,
		logger: o.logger.Named("ip sync").With(
			zap.Int("period", o.period),
		),
		period: o.period,
	}
}

func (w *syncer) Name() string {
	return "ip"
}

func (w *syncer) Run(ctx context.Context) error {
	ticker := time.NewTicker(time.Duration(w.period) * time.Second)
	w.logger.Info("start ip sync")

	for {
		select {
		case <-ticker.C:
			w.logger.Info("triggered ip sync")
			if err := w.sync(ctx); err != nil {
				w.logger.Error("failed to sync ip", zap.Error(err))
			}
		case <-ctx.Done():
			w.logger.Info("finished ip sync", zap.Error(ctx.Err()))
			return nil
		}
	}
}

func (w *syncer) sync(ctx context.Context) error {
	return w.dmu.SyncNodes(ctx)
}
