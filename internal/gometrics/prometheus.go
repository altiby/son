package metrics

import (
	"strings"

	met "github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

const (
	// AppNameLabel - standard label for application name.
	AppNameLabel = "appName"
	// HTTPCodeLabel - label for HTTP status code
	HTTPCodeLabel = "code"
	// HTTPMethodLabel - label for HTTP method (GET, POST, etc.)
	HTTPMethodLabel = "method"
	// HandlerLabel - label to differentiate handlers, this can be router path, or handler name.
	HandlerLabel = "handler"
)

// NewPrometheusDriver – constructor for Prometheus driver.
// ExponentialBuckets(0.25, 4, 30) have been chosen as a compromise for request latency measuring
// if other bucket set is required, it could be passed to Histogram() or Elapsed() call
func NewPrometheusDriver(constLabels map[string]string) Driver {
	return &prometheusDriver{
		constLabels:    stdprometheus.Labels(constLabels),
		defaultBuckets: stdprometheus.ExponentialBuckets(0.25, 4, 30),
	}
}

// NewPrometheusMetrics – provide prometheus implementation for collecting metrics.
func NewPrometheusMetrics(appName string) Metrics {
	driver := NewPrometheusDriver(map[string]string{AppNameLabel: appName})
	m := NewMetrics(driver)
	ExportRunningGoRoutines(m)
	return m
}

type prometheusDriver struct {
	constLabels    stdprometheus.Labels
	defaultBuckets []float64
}

// NewCounter – constructor for counter collector.
func (p *prometheusDriver) NewCounter(metric string, tags []string, sampleRate float64) met.Counter {
	return prometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Name:        strings.Replace(metric, ".", ":", -1),
		Help:        strings.Replace(metric, ".", ":", -1),
		ConstLabels: p.constLabels,
	}, tagsIntoTagNames(tags))
}

// NewHistogram – constructor for histogram collector.
func (p *prometheusDriver) NewHistogram(metric string, tags []string, sampleRate float64, buckets ...[]float64) met.Histogram {
	var histogramBuckets = p.defaultBuckets

	// a bit ugly way to pass an optional argument and not to break the existing code
	if len(buckets) > 0 {
		histogramBuckets = buckets[0]
	}
	return prometheus.NewHistogramFrom(stdprometheus.HistogramOpts{
		Name:        strings.Replace(metric, ".", ":", -1),
		Help:        strings.Replace(metric, ".", ":", -1),
		ConstLabels: p.constLabels,
		Buckets:     histogramBuckets,
	}, tagsIntoTagNames(tags))
}

// NewGauge – constructor for gauge collector.
func (p *prometheusDriver) NewGauge(metric string, tags []string) met.Gauge {
	return prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
		Name:        strings.Replace(metric, ".", ":", -1),
		Help:        strings.Replace(metric, ".", ":", -1),
		ConstLabels: p.constLabels,
	}, tagsIntoTagNames(tags))
}
