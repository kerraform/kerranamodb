package dlock

import (
	"context"
	"net"

	"github.com/minio/dsync/v3"
)

type options struct {
	sd        string
	endpoints []string
}

type Lock struct {
	endpoints []string
	lock      *dsync.Dsync
	ips       []string
}

type LockOptions func(*options)

func WithServiceDiscovery(sd string) LockOptions {
	return func(o *options) {
		o.sd = sd
	}
}

func WithStaticIPs(endpoints []string) LockOptions {
	return func(o *options) {
		o.endpoints = endpoints
	}
}

func NewLock(ctx context.Context, opts ...LockOptions) (*Lock, error) {
	var o options
	for _, opt := range opts {
		opt(&o)
	}

	var err error
	var eps []string
	var lks []dsync.NetLocker
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

	ds, err := dsync.New(lks)
	if err != nil {
		return nil, err
	}

	return &Lock{
		endpoints: eps,
		lock:      ds,
	}, nil
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
