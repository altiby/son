package metrics

import (
	"fmt"
	"sort"
	"strings"
	"time"

	met "github.com/go-kit/kit/metrics"
)

// NewMetrics construct metrics builder with included metrics driver.
func NewMetrics(driver Driver) *MC {
	return &MC{
		driver:     driver,
		gauges:     make(map[string]met.Gauge),
		counters:   make(map[string]met.Counter),
		histograms: make(map[string]met.Histogram),
	}
}

// MC – metrics controller.
type MC struct {
	driver     Driver
	counters   map[string]met.Counter
	histograms map[string]met.Histogram
	gauges     map[string]met.Gauge
}

func getMetricName(m string, t []string) string {
	sort.Strings(t)
	return fmt.Sprintf("%s|%s", m, strings.Join(t, ","))
}

// Timing ...
func Timing(metricBase string) string {
	return metricBase + ".seconds"
}

// Counter – update counter value.
func (mc *MC) Counter(metric string, tags []string, count float64) {
	storedMetric := getMetricName(metric, tagsIntoTagNames(tags))
	if mc.counters[storedMetric] == nil {
		counter := mc.driver.NewCounter(metric, tags, 1)
		mc.counters[storedMetric] = counter
	}
	mc.counters[storedMetric].With(tags...).Add(count)
}

// Elapsed reports the time elapsed for a function call in nanoseconds
func (mc *MC) Elapsed(metric string, tags []string, buckets ...[]float64) func(time.Time) {
	return func(begin time.Time) {
		// use base units: https://prometheus.io/docs/practices/naming/#base-units
		seconds := float64(time.Since(begin).Nanoseconds()) / float64(time.Second)
		mc.Histogram(Timing(metric), tags, seconds, buckets...)
	}
}

// Histogram measure the statistical distribution of a set of values
func (mc *MC) Histogram(metric string, tags []string, value float64, buckets ...[]float64) {
	storedMetric := getMetricName(metric, tagsIntoTagNames(tags))
	if mc.histograms[storedMetric] == nil {
		histogram := mc.driver.NewHistogram(metric, tags, 1, buckets...)
		mc.histograms[storedMetric] = histogram
	}
	mc.histograms[storedMetric].With(tags...).Observe(value)
}

func (mc *MC) getOrCreateGauge(metric string, tags []string) met.Gauge {
	storedMetric := getMetricName(metric, tagsIntoTagNames(tags))
	if mc.gauges[storedMetric] == nil {
		gauge := mc.driver.NewGauge(metric, tags)
		mc.gauges[storedMetric] = gauge
	}
	return mc.gauges[storedMetric].With(tags...)
}

// GaugeSet set gauge value.
func (mc *MC) GaugeSet(metric string, tags []string, value float64) {
	mc.getOrCreateGauge(metric, tags).Set(value)
}

// GaugeAdd update gauge value.
func (mc *MC) GaugeAdd(metric string, tags []string, value float64) {
	mc.getOrCreateGauge(metric, tags).Add(value)
}

func tagsIntoTagNames(tags []string) []string {
	var tagNames []string
	for i := 0; i < len(tags); i++ {
		if i%2 == 0 {
			tagNames = append(tagNames, tags[i])
		}
	}
	return tagNames
}
