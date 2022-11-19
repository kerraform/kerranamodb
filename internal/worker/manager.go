package worker

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type Manager struct {
	workers []Worker
}

func NewManager(ws ...Worker) *Manager {
	return &Manager{
		workers: ws,
	}
}

func (m *Manager) Run(ctx context.Context) error {
	wg, ctx := errgroup.WithContext(ctx)

	for _, w := range m.workers {
		w := w
		wg.Go(func() error {
			return w.Run(ctx)
		})
	}

	return wg.Wait()
}

// Names return the list of workers registered
func (m *Manager) Names() []string {
	res := make([]string, len(m.workers))
	for i, w := range m.workers {
		res[i] = w.Name()
	}

	return res
}
