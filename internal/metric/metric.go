package metric

import (
	"context"
	"fmt"
	"time"

	"github.com/kerraform/kerranamodb/internal/driver"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type MetricName string
type MetricLabel string

const (
	metricNamespace = "kegistry"

	// Metrics
	metricNameHTTPRequestTotal MetricName = "registry_request_total"

	// Labels
	metricLabelHTTPStatusCode MetricLabel = "code"
	metricLabelHTTPMethod     MetricLabel = "method"
	metricLabelHTTPPath       MetricLabel = "path"
)

type MetricSyncFunc func(m *RegistryMetrics)

type RegistryMetrics struct {
	driver  driver.Driver
	logger  *zap.Logger
	metrics map[MetricName]prometheus.Collector
}

func New(logger *zap.Logger, driver driver.Driver) *RegistryMetrics {
	return &RegistryMetrics{
		driver: driver,
		logger: logger.Named("metric"),
		metrics: map[MetricName]prometheus.Collector{
			metricNameHTTPRequestTotal: prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Namespace: metricNamespace,
					Name:      string(metricNameHTTPRequestTotal),
					Help:      "Total count of the request",
				},
				[]string{
					string(metricLabelHTTPStatusCode),
					string(metricLabelHTTPMethod),
					string(metricLabelHTTPPath),
				},
			),
		},
	}
}

// RegisterAllMetrics records the current number of registries
func (m *RegistryMetrics) RegisterAllMetrics() {
	for _, pm := range m.metrics {
		prometheus.MustRegister(pm)
	}
}

// IncrementHTTPRequestTotal increment
func (m *RegistryMetrics) IncrementHTTPRequestTotal(code int, method, path string) {
	if c, ok := m.metrics[metricNameHTTPRequestTotal].(*prometheus.CounterVec); ok {
		c.WithLabelValues(fmt.Sprint(code), method, path).Add(1)
	}
}

// Resync recomputes the metrics
func (m *RegistryMetrics) Resync(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				if err := m.resync(ctx); err != nil {
					m.logger.Error("error collecting prometheus metrics", zap.Error(err))
				}
				m.logger.Debug("collected prometheus metrics")
			case <-ctx.Done():
			}
		}
	}()
}

func (m *RegistryMetrics) resync(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	wg, ctx := errgroup.WithContext(ctx)
	return wg.Wait()
}
