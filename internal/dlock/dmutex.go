package dlock

import (
	"sync"

	"github.com/minio/dsync/v3"
	"go.uber.org/zap"
)

type dmutex struct {
	dmu       *dsync.DRWMutex
	logger    *zap.Logger
	mu        *sync.RWMutex
	isReading bool
	isWriting bool
}

func (d *dmutex) Lock() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.isWriting = true
	d.logger.Debug("lock", zap.Bool("isWriting", d.isWriting), zap.Bool("isReading", d.isReading))
}

func (d *dmutex) Unlock() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.isWriting = false
	d.logger.Debug("unlock", zap.Bool("isWriting", d.isWriting), zap.Bool("isReading", d.isReading))
}

func (d *dmutex) Rlock() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.isReading = true
	d.logger.Debug("rlock", zap.Bool("isWriting", d.isWriting), zap.Bool("isReading", d.isReading))
}

func (d *dmutex) RUnlock() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.isReading = false
	d.logger.Debug("runlock", zap.Bool("isWriting", d.isWriting), zap.Bool("isReading", d.isReading))
}
