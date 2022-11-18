package dlock

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/minio/dsync/v3"
	"go.uber.org/zap"
)

type options struct {
	sd        string
	endpoints []string
	logger    *zap.Logger
}

type DMutex struct {
	DSync *dsync.Dsync
	Ready bool

	endpoints []string
	logger    *zap.Logger
}

type LockOptions func(*options)

func WithLogger(logger *zap.Logger) LockOptions {
	return func(o *options) {
		o.logger = logger
	}
}

func WithServiceDiscovery(sd string) LockOptions {
	return func(o *options) {
		o.sd = sd
	}
}

func WithStaticEndpoints(endpoints []string) LockOptions {
	return func(o *options) {
		o.endpoints = endpoints
	}
}

func NewDMutex(ctx context.Context, opts ...LockOptions) (*DMutex, error) {
	o := &options{
		logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(o)
	}

	var eps []string
	if o.sd != "" {
		ips, err := fetchNodes(o.sd)
		if err != nil {
			return nil, err
		}
		for _, ip := range ips {
			eps = append(eps, ip.String())
		}
	}

	if len(o.endpoints) > 0 {
		eps = append(eps, o.endpoints...)
	}

	return &DMutex{
		endpoints: eps,
		logger:    o.logger,
	}, nil
}

func (dmu *DMutex) Connect(ctx context.Context) error {
	var lks []dsync.NetLocker

	var wg sync.WaitGroup
	for _, e := range dmu.endpoints {
		wg.Add(1)
		e := e
		go func() {
			l, err := dmu.connect(ctx, &DLockerConfig{
				endpoint: e,
				logger:   dmu.logger.Named("dlocker").With(zap.String("endpoint", e)),
			})
			if err != nil {
				return
			}

			lks = append(lks, l)
			wg.Done()
		}()
	}

	wg.Wait()
	ds, err := dsync.New(lks)
	if err != nil {
		dmu.logger.Debug("failed to connect", zap.Error(err))
		return err
	}

	dmu.logger.Info("dsync initialized", zap.Int("node", len(dmu.endpoints)))
	dmu.DSync = ds
	dmu.Ready = true
	return nil
}

func (dmu *DMutex) connect(ctx context.Context, cfg *DLockerConfig) (dsync.NetLocker, error) {
	l := NewDLocker(ctx, cfg)
	dmu.logger.Info("connected to node", zap.String("endpoint", cfg.endpoint))
	return l, nil
}

func (dmu *DMutex) Lock(ctx context.Context, dlid DLockID) (*dsync.DRWMutex, error) {
	mu := dsync.NewDRWMutex(ctx, string(dlid), dmu.DSync)
	ch := make(chan bool)
	defer close(ch)

	go func(id DLockID) {
		ch <- mu.GetLock(string(dlid), "", 1*time.Second)
	}(dlid)

	select {
	case success := <-ch:
		if success {
			return mu, nil
		}

		return nil, fmt.Errorf("state is locked")
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (dmu *DMutex) RLock(ctx context.Context, dlid DLockID) (*dsync.DRWMutex, error) {
	mu := dsync.NewDRWMutex(ctx, string(dlid), dmu.DSync)
	ch := make(chan bool)
	defer close(ch)

	go func(id DLockID) {
		ch <- mu.GetRLock(string(dlid), "", 1*time.Second)
	}(dlid)

	select {
	case success := <-ch:
		if success {
			return mu, nil
		}

		return nil, fmt.Errorf("state is locked")
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func fetchNodes(endpoint string) ([]net.IP, error) {
	res := []net.IP{}

	ips, err := net.LookupIP(endpoint)
	if err != nil {
		return nil, err
	}

	for _, ip := range ips {
		if ipv4 := ip.To4(); ipv4 != nil {
			res = append(res, ipv4)
		}
	}

	return res, nil
}
