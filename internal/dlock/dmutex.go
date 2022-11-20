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
	d.logger.Debug("lock", zap.Bool("isWriting", d.isWriting), zap.Bool("isReading", d.isReading))
	d.mu.Lock()
	d.isWriting = true
}

func (d *dmutex) Unlock() {
	d.logger.Debug("unlock", zap.Bool("isWriting", d.isWriting), zap.Bool("isReading", d.isReading))
	d.isWriting = false
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
