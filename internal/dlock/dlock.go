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

	mu        *sync.RWMutex
	mus       map[DLockID]*dmutex
	endpoints []string
	logger    *zap.Logger
}

type dmutex struct {
	mu        *sync.RWMutex
	isReading bool
	isWriting bool
}

func (d *dmutex) Lock() {
	d.mu.Lock()
	d.isWriting = true
}

func (d *dmutex) UnLock() {
	d.isWriting = true
	d.mu.Unlock()
}

func (d *dmutex) Rlock() {
	d.mu.RLock()
	d.isReading = true
}

func (d *dmutex) RUnlock() {
	d.isReading = false
	d.mu.RUnlock()
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
		mu:        &sync.RWMutex{},
		mus:       map[DLockID]*dmutex{},
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

func (dmu *DMutex) SetReading(lid DLockID) {
	dmu.logger.Debug("set to reading status")
	dmu.setReading(lid, true)
}

func (dmu *DMutex) SetUnReading(lid DLockID) {
	dmu.logger.Debug("set to unreading status")
	dmu.setReading(lid, false)
}

func (dmu *DMutex) setReading(lid DLockID, v bool) {
	dmu.mu.Lock()
	defer func() {
		dmu.mu.Unlock()
		dmu.logger.Debug("unlock root mutex")
	}()
	// mu, ok := dmu.mus[lid]
	// if ok {
	// if v {
	// 	dmu.logger.Debug("rlocked", zap.String("dlid", string(lid)))
	// 	mu.Rlock()
	// 	return
	// }

	// dmu.logger.Debug("runlocked", zap.String("dlid", string(lid)))
	// mu.RUnlock()
	// return
	// }

	// dmu.mus[lid] = &dmutex{
	// 	mu:        &sync.RWMutex{},
	// 	isReading: v,
	// }

	dmu.logger.Debug("initial set reading", zap.String("dlid", string(lid)))
}

func (dmu *DMutex) IsReadable(lid DLockID) bool {
	dmu.mu.RLock()
	defer dmu.mu.RUnlock()
	mu, ok := dmu.mus[lid]
	if !ok {
		return false
	}

	mu.mu.RLock()
	defer mu.mu.RUnlock()
	return mu.isReading || mu.isWriting
}

func (dmu *DMutex) SetWriting(lid DLockID) {
	dmu.setWriting(lid, true)
}

func (dmu *DMutex) SetUnWriting(lid DLockID) {
	dmu.setWriting(lid, false)
}

func (dmu *DMutex) setWriting(lid DLockID, v bool) {
	dmu.mu.Lock()
	defer dmu.mu.Unlock()
	mu, ok := dmu.mus[lid]
	if ok {
		if v {
			mu.Lock()
			return
		}

		mu.UnLock()
		return
	}

	dmu.mus[lid] = &dmutex{
		mu:        &sync.RWMutex{},
		isWriting: v,
		isReading: v,
	}

	dmu.logger.Debug("initial write reading", zap.String("dlid", string(lid)))
}

func (dmu *DMutex) IsWritable(lid DLockID) bool {
	dmu.mu.RLock()
	defer dmu.mu.RUnlock()
	mu, ok := dmu.mus[lid]
	if !ok {
		return false
	}

	mu.mu.RLock()
	defer mu.mu.RUnlock()
	return mu.isWriting
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
			dmu.logger.Debug("success to get lock", zap.Bool("success", success))
			return mu, nil
		}

		return nil, fmt.Errorf("state is locked")
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (dmu *DMutex) RLock(ctx context.Context, dlid DLockID) (*dsync.DRWMutex, error) {
	dmu.SetReading(dlid)
	defer dmu.SetUnReading(dlid)

	mu := dsync.NewDRWMutex(ctx, string(dlid), dmu.DSync)
	ch := make(chan bool)
	defer close(ch)

	go func(dlid DLockID) {
		ch <- mu.GetRLock(string(dlid), "", 1*time.Second)
	}(dlid)

	select {
	case success := <-ch:
		if success {
			dmu.logger.Debug("success to get rlock", zap.Bool("success", success))
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
