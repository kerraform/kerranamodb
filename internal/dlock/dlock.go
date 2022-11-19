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
	expected  int
	hostIP    string
	sd        string
	port      int
	endpoints []string
	logger    *zap.Logger
	timeout   *time.Duration
}

type DMutex struct {
	DSync *dsync.Dsync
	Ready bool

	mu  *sync.RWMutex
	mus map[DLockID]*dmutex

	endpoints []string
	expected  int
	logger    *zap.Logger
	hostIP    string
}

type dmutex struct {
	logger    *zap.Logger
	mu        *sync.RWMutex
	isReading bool
	isWriting bool
}

func (d *dmutex) Lock() {
	d.logger.Debug("lock", zap.Bool("isWriting", d.isWriting), zap.Bool("isReading", d.isReading))
	d.mu.Lock()
	d.isWriting = true
}

func (d *dmutex) UnLock() {
	d.logger.Debug("unlock", zap.Bool("isWriting", d.isWriting), zap.Bool("isReading", d.isReading))
	d.isWriting = true
	d.mu.Unlock()
}

func (d *dmutex) Rlock() {
	d.logger.Debug("rlock", zap.Bool("isWriting", d.isWriting), zap.Bool("isReading", d.isReading))
	d.mu.RLock()
	d.isReading = true
}

func (d *dmutex) RUnlock() {
	d.logger.Debug("runlock", zap.Bool("isWriting", d.isWriting), zap.Bool("isReading", d.isReading))
	d.isReading = false
	d.mu.RUnlock()
}

type LockOptions func(*options)

func WithLogger(logger *zap.Logger) LockOptions {
	return func(o *options) {
		o.logger = logger
	}
}

func WithServiceDiscovery(sd string, count int, hostIP string, port int) LockOptions {
	return func(o *options) {
		o.expected = count
		o.port = port
		o.hostIP = hostIP
		o.sd = sd
	}
}

func WithStaticEndpoints(endpoints []string) LockOptions {
	return func(o *options) {
		o.endpoints = endpoints
	}
}

func WithTimeout(timeout time.Duration) LockOptions {
	return func(o *options) {
		o.timeout = &timeout
	}
}

func NewDMutex(ctx context.Context, opts ...LockOptions) (*DMutex, error) {
	o := &options{
		logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(o)
	}

	if o.timeout != nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, *o.timeout)
		defer cancel()
	}

	dmu := &DMutex{
		expected: o.expected,
		logger:   o.logger,
		mu:       &sync.RWMutex{},
		mus:      map[DLockID]*dmutex{},
	}

	var eps []string
	if o.sd != "" {
		ips, err := dmu.fetchNodes(ctx, o.sd)
		if err != nil {
			return nil, err
		}
		for _, ip := range ips {
			eps = append(eps, fmt.Sprintf("http://%s:%d", ip.String(), o.port))
		}
	}

	if len(o.endpoints) > 0 {
		eps = append(eps, o.endpoints...)
	}

	dmu.endpoints = eps
	return dmu, nil
}

func (dmu *DMutex) Connect(ctx context.Context) error {
	lks := []dsync.NetLocker{}

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
	defer dmu.mu.Unlock()
	mu, ok := dmu.mus[lid]
	if ok {
		if v {
			mu.Rlock()
			return
		}

		mu.RUnlock()
		return
	}

	m := &dmutex{
		logger:    dmu.logger.Named("dmutex").With(zap.String("dlid", string(lid))),
		mu:        &sync.RWMutex{},
		isReading: v,
	}

	if v {
		m.Rlock()
	} else {
		m.RUnlock()
	}

	dmu.mus[lid] = m
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

	m := &dmutex{
		logger:    dmu.logger.Named("dmutex").With(zap.String("dlid", string(lid))),
		mu:        &sync.RWMutex{},
		isWriting: v,
		isReading: v,
	}

	if v {
		m.Rlock()
	} else {
		m.RUnlock()
	}

	dmu.mus[lid] = m
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
	dmu.SetWriting(dlid)
	defer dmu.SetUnWriting(dlid)

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

func (dmu *DMutex) fetchNodes(ctx context.Context, endpoint string) ([]net.IP, error) {
	res := []net.IP{}

	dmu.logger.Debug("fetching service discovery endpoint", zap.String("endpoint", endpoint))

LOOP:
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(10000 * time.Millisecond):
			ips, err := net.LookupIP(endpoint)
			if err != nil {
				dmu.logger.Warn("failed to lookup ip", zap.String("endpoint", endpoint), zap.Error(err))
				continue
			}

			dmu.logger.Debug("service discovery ips", zap.Int("count", len(ips)), zap.Int("expected", dmu.expected))
			if len(ips) == dmu.expected {
				for _, ip := range ips {
					if ipv4 := ip.To4(); ipv4 != nil {
						dmu.logger.Debug("fetched ip", zap.String("ip", ipv4.String()))
						if ipv4.String() != dmu.hostIP {
							res = append(res, ipv4)
						}
					}
				}

				break LOOP
			}
		}
	}

	return res, nil
}
