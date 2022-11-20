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
	port      int
	logger    *zap.Logger
	sd        string
	hostIP    string
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
		hostIP:   o.hostIP,
		mu:       &sync.RWMutex{},
		mus:      map[DLockID]*dmutex{},
		port:     o.port,
		sd:       o.sd,
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

func (dmu *DMutex) SetReading(ctx context.Context, dlid DLockID) *dmutex {
	dmu.logger.Debug("set to reading status")
	return dmu.setReading(ctx, dlid, true)
}

func (dmu *DMutex) SetUnReading(ctx context.Context, dlid DLockID) *dmutex {
	dmu.logger.Debug("set to unreading status")
	return dmu.setReading(ctx, dlid, false)
}

func (dmu *DMutex) setReading(ctx context.Context, dlid DLockID, lock bool) *dmutex {
	dmu.mu.Lock()
	defer dmu.mu.Unlock()
	mu, ok := dmu.mus[dlid]
	if !ok {
		mu = &dmutex{
			dmu:       dsync.NewDRWMutex(ctx, string(dlid), dmu.DSync),
			logger:    dmu.logger.Named("dmutex").With(zap.String("dlid", string(dlid))),
			mu:        &sync.RWMutex{},
			isReading: lock,
		}
	}

	if lock {
		mu.Rlock()
	} else {
		mu.RUnlock()
	}

	dmu.mus[dlid] = mu
	dmu.logger.Debug("initial set reading", zap.String("dlid", string(dlid)))
	return mu
}

func (dmu *DMutex) IsReadable(lid DLockID) bool {
	dmu.mu.RLock()
	defer dmu.mu.RUnlock()
	mu, ok := dmu.mus[lid]
	if !ok {
		return true
	}

	mu.mu.RLock()
	defer mu.mu.RUnlock()
	return !(mu.isReading || mu.isWriting)
}

func (dmu *DMutex) SetWriting(ctx context.Context, lid DLockID) *dmutex {
	return dmu.setWriting(ctx, lid, true)
}

func (dmu *DMutex) SetUnWriting(ctx context.Context, lid DLockID) *dmutex {
	return dmu.setWriting(ctx, lid, false)
}

func (dmu *DMutex) setWriting(ctx context.Context, dlid DLockID, lock bool) *dmutex {
	dmu.mu.Lock()
	defer dmu.mu.Unlock()
	mu, ok := dmu.mus[dlid]
	if !ok {
		mu = &dmutex{
			dmu:       dsync.NewDRWMutex(ctx, string(dlid), dmu.DSync),
			logger:    dmu.logger.Named("dmutex").With(zap.String("dlid", string(dlid))),
			mu:        &sync.RWMutex{},
			isWriting: lock,
			isReading: lock,
		}
	}

	if lock {
		mu.Lock()
	} else {
		mu.Unlock()
	}

	dmu.mus[dlid] = mu
	dmu.logger.Debug("initial set write", zap.String("dlid", string(dlid)))
	return mu
}

func (dmu *DMutex) IsWritable(lid DLockID) bool {
	dmu.mu.RLock()
	defer dmu.mu.RUnlock()
	mu, ok := dmu.mus[lid]
	if !ok {
		return true
	}

	mu.mu.RLock()
	defer mu.mu.RUnlock()
	return !mu.isWriting
}

func (dmu *DMutex) connect(ctx context.Context, cfg *DLockerConfig) (dsync.NetLocker, error) {
	l := NewDLocker(ctx, cfg)
	dmu.logger.Info("connected to node", zap.String("endpoint", cfg.endpoint))
	return l, nil
}

func (dmu *DMutex) Lock(ctx context.Context, dlid DLockID) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	d := dmu.SetWriting(ctx, dlid)
	defer dmu.SetUnWriting(ctx, dlid)

	ch := make(chan bool)
	defer close(ch)

	go func(id DLockID) {
		ch <- d.dmu.GetLock(string(dlid), "", 1*time.Second)
	}(dlid)

	select {
	case success := <-ch:
		if success {
			dmu.logger.Debug("acquired lock", zap.Bool("success", success))
			return nil
		}

		return fmt.Errorf("state is locked")
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (dmu *DMutex) Unlock(ctx context.Context, dlid DLockID) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	d := dmu.SetWriting(ctx, dlid)
	defer dmu.SetUnWriting(ctx, dlid)

	ch := make(chan struct{})
	defer close(ch)

	go func() {
		d.dmu.Unlock()
		select {
		case <-ch:
		default:
			ch <- struct{}{}
		}
	}()

	select {
	case <-ch:
		dmu.logger.Debug("released lock")
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// RLock read locks for the passed dlock ID and writes the status
// to "isReading".
func (dmu *DMutex) RLock(ctx context.Context, dlid DLockID) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	d := dmu.SetReading(ctx, dlid)
	defer dmu.SetUnReading(ctx, dlid)

	ch := make(chan bool)
	defer close(ch)

	go func(dlid DLockID) {
		ch <- d.dmu.GetRLock(string(dlid), "", 1*time.Second)
	}(dlid)

	select {
	case success := <-ch:
		if success {
			dmu.logger.Debug("acquired read lock", zap.Bool("success", success))
			return nil
		}

		return fmt.Errorf("state is locked")
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (dmu *DMutex) RUnlock(ctx context.Context, dlid DLockID) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	d := dmu.SetReading(ctx, dlid)
	defer dmu.SetUnReading(ctx, dlid)

	ch := make(chan struct{})
	defer close(ch)

	go func() {
		d.dmu.RUnlock()
		ch <- struct{}{}
	}()

	select {
	case <-ch:
		dmu.logger.Debug("released read lock")
		return nil
	case <-ctx.Done():
		return ctx.Err()
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

func (dmu *DMutex) SyncNodes(ctx context.Context) error {
	ips, err := net.LookupIP(dmu.sd)
	if err != nil {
		return err
	}

	dmu.logger.Debug("service discovery ips", zap.Int("count", len(ips)), zap.Int("expected", dmu.expected))
	var eps []string
	if len(ips) == dmu.expected {
		for _, ip := range ips {
			if ipv4 := ip.To4(); ipv4 != nil {
				dmu.logger.Debug("fetched ip", zap.String("ip", ipv4.String()))
				if ipv4.String() != dmu.hostIP {
					eps = append(eps, fmt.Sprintf("http://%s:%d", ip.String(), dmu.port))
				}
			}
		}
	}

	lks := []dsync.NetLocker{}
	var wg sync.WaitGroup
	for _, e := range eps {
		wg.Add(1)
		e := e
		go func() {
			l, err := dmu.connect(ctx, &DLockerConfig{
				endpoint: e,
				logger:   dmu.logger.Named("dlocker").With(zap.String("endpoint", e)),
			})
			if err != nil {
				dmu.logger.Warn("failed to connect to lock node", zap.Error(err))
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

	dmu.mu.Lock()
	defer dmu.mu.Unlock()

	dmu.logger.Info("updated lock nodes", zap.Int("nodes", len(eps)))
	dmu.endpoints = eps
	dmu.DSync = ds
	return nil
}
