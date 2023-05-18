package metrics

import (
	"time"

	met "github.com/go-kit/kit/metrics"
)

// Metrics – functional to collect metrics values.
type Metrics interface {
	Elapsed(metric string, tags []string, buckets ...[]float64) func(time.Time)
	Counter(metric string, tags []string, value float64)
	Histogram(metric string, tags []string, value float64, buckets ...[]float64)
	GaugeSet(metric string, tags []string, value float64)
	GaugeAdd(metric string, tags []string, value float64)
}

// Driver – functional to construct values' collectors.
type Driver interface {
	NewCounter(string, []string, float64) met.Counter
	NewHistogram(string, []string, float64, ...[]float64) met.Histogram
	NewGauge(string, []string) met.Gauge
}
