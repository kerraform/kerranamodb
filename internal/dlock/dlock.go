package dlock

import (
	"context"
	"net"
	"sync"

	"github.com/minio/dsync/v3"
	"go.uber.org/zap"
)

type options struct {
	sd        string
	endpoints []string
	logger    *zap.Logger
}

type DMutex struct {
	endpoints []string
Lock      *dsync.Dsync
	ips       []string
	ready     bool
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
			l, err := dmu.connect(ctx, e)
			if err != nil {
				return
			}

			lks = append(lks, l)
			wg.Done()
		}()
	}

	ds, err := dsync.New(lks)
	if err != nil {
		return nil
	}

	dmu.Lock = ds
	dmu.ready = true
	return nil
}

func (dmu *DMutex) connect(ctx context.Context, endpoint string) (dsync.NetLocker, error) {
	l := NewDLocker(ctx, endpoint)
	return l, nil
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
