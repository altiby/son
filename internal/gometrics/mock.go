package metrics

import (
	"time"

	met "github.com/go-kit/kit/metrics"
)

// MockMetrics – mock structure to collect metrics.
type MockMetrics struct {
	Metrics map[string]float64
}

// NewDefaultMockMetrics constructor for MockMetrics with default config values.
func NewDefaultMockMetrics() Metrics {
	m := NewMockMetrics()
	return m
}

// NewMockMetrics constructor for MockMetrics.
func NewMockMetrics() *MockMetrics {
	return &MockMetrics{Metrics: make(map[string]float64)}
}

// Counter – update counter value.
func (m *MockMetrics) Counter(metric string, tags []string, count float64) {
	m.Metrics[getMetricName(metric, tags)] = count
}

// GaugeAdd – update counter value.
func (m *MockMetrics) GaugeAdd(metric string, tags []string, v float64) {
	m.Metrics[getMetricName(metric, tags)] = v
}

// GaugeSet – update counter value.
func (m *MockMetrics) GaugeSet(metric string, tags []string, v float64) {
	m.Metrics[getMetricName(metric, tags)] = v
}

// Elapsed – update elapsed value.
func (m *MockMetrics) Elapsed(metric string, tags []string, buckets ...[]float64) func(time.Time) {
	return func(begin time.Time) {
	}
}

// Histogram – update elapsed value.
func (m *MockMetrics) Histogram(metric string, tags []string, value float64, buckets ...[]float64) {
	m.Metrics[getMetricName(metric, tags)] = value
}

// MockDriver implements mock functional to construct values' collectors and publish them.
type MockDriver struct {
	NewCounterFn   func(string, float64) met.Counter
	NewHistogramFn func(string, float64) met.Histogram
	NewGaugeFn     func(string) met.Gauge

	MetricName string
}

// NewCounter – constructor for mock counter collector.
func (d *MockDriver) NewCounter(metric string, sampleRate float64) met.Counter {
	d.MetricName = metric
	return d.NewCounterFn(metric, sampleRate)
}

// NewHistogram – constructor for mock histogram collector.
func (d *MockDriver) NewHistogram(metric string, sampleRate float64, buckets ...[]float64) met.Histogram {
	d.MetricName = metric
	return d.NewHistogramFn(metric, sampleRate)
}

// NewGauge – constructor for mock gauge collector.
func (d *MockDriver) NewGauge(metric string) met.Gauge {
	d.MetricName = metric
	return d.NewGaugeFn(metric)
}

// MockCounter – mock counter collector.
type MockCounter struct {
	WithFn func(labelValues ...string) met.Counter
	AddFn  func(delta float64)

	Labels []string
	Value  float64
}

// With function create contextual counter collector.
func (c *MockCounter) With(labelValues ...string) met.Counter {
	return c.WithFn(labelValues...)
}

// Add function update counter value.
func (c *MockCounter) Add(delta float64) {
	c.AddFn(delta)
}

// MockGauge – mock gauge collector.
type MockGauge struct {
	WithFn func(labelValues ...string) met.Gauge
	SetFn  func(value float64)
	AddFn  func(delta float64)

	Labels   []string
	SetValue float64
	AddValue float64
}

// With function create contextual gauge collector.
func (g *MockGauge) With(labelValues ...string) met.Gauge {
	return g.WithFn(labelValues...)
}

// Set gauge value.
func (g *MockGauge) Set(value float64) {
	g.SetFn(value)
}

// Add function update gauge value.
func (g *MockGauge) Add(delta float64) {
	g.AddFn(delta)
}

// MockHistogram – mock histogram collector.
type MockHistogram struct {
	WithFn    func(labelValues ...string) met.Histogram
	ObserveFn func(value float64)

	Labels []string
	Value  float64
}

// With function create contextual histogram collector.
func (h *MockHistogram) With(labelValues ...string) met.Histogram {
	return h.WithFn(labelValues...)
}

// Observe histogram value.
func (h *MockHistogram) Observe(value float64) {
	h.ObserveFn(value)
}
