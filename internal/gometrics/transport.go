package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	inFlightRequests       = "in_flight_requests"
	apiRequestsTotal       = "requests_total"
	requestDurationSeconds = "request_duration_seconds"
	responseSizeBytes      = "response_size_bytes"
)

// HandlerMetrics models the prometheus metrics for HTTP handlers
type HandlerMetrics struct {
	appName       string
	inFlightGauge prometheus.Gauge
	counter       *prometheus.CounterVec
	duration      *prometheus.HistogramVec
	responseSize  *prometheus.HistogramVec
}

// NewHandlerMetrics initializes the HandlerMetrics
func NewHandlerMetrics(appName string) HandlerMetrics {
	hm := HandlerMetrics{
		appName: appName,
		inFlightGauge: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: inFlightRequests,
				Help: "A gauge of requests currently being served by the wrapped handler.",
				ConstLabels: prometheus.Labels{
					AppNameLabel: appName,
				},
			},
		),
		counter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: apiRequestsTotal,
				Help: "A counter for requests to the wrapped handler.",
				ConstLabels: prometheus.Labels{
					AppNameLabel: appName,
				},
			},
			[]string{HTTPCodeLabel, HTTPMethodLabel, HandlerLabel},
		),
		// duration is partitioned by the HTTP method and handler. It uses custom buckets based on the expected request duration.
		duration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: requestDurationSeconds,
				Help: "A histogram of latencies for requests.",
				ConstLabels: prometheus.Labels{
					AppNameLabel: appName,
				},
			},
			[]string{HandlerLabel, HTTPMethodLabel},
		),
		// responseSize has no labels, making it a zero-dimensional ObserverVec.
		responseSize: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: responseSizeBytes,
				Help: "A histogram of response sizes for requests.",
				ConstLabels: prometheus.Labels{
					AppNameLabel: appName,
				},
			},
			[]string{HandlerLabel, HTTPMethodLabel},
		),
	}
	// Ensuring the collectors to be registered only once
	ensureRegister(hm.inFlightGauge, hm.counter, hm.duration, hm.responseSize)
	return hm
}

// Handler instruments the handlers with all the metrics, injecting the "handler" label by currying.
func (hm HandlerMetrics) Handler(method string, path string, h http.Handler) http.Handler {
	h = promhttp.InstrumentHandlerResponseSize(hm.responseSize.MustCurryWith(prometheus.Labels{
		HandlerLabel:    path,
		HTTPMethodLabel: method,
	}), h)

	h = promhttp.InstrumentHandlerCounter(hm.counter.MustCurryWith(prometheus.Labels{
		HandlerLabel:    path,
		HTTPMethodLabel: method,
	}), h)

	h = promhttp.InstrumentHandlerDuration(hm.duration.MustCurryWith(prometheus.Labels{
		HandlerLabel:    path,
		HTTPMethodLabel: method,
	}), h)

	return promhttp.InstrumentHandlerInFlight(hm.inFlightGauge, h)
}

// ensureRegister is a less strict implementation of prometheus.MustRegister as it doesn't panic
// in case the error returned by the Register function is an AlreadyRegisteredError instance.
func ensureRegister(collectors ...prometheus.Collector) {
	for _, c := range collectors {
		if err := prometheus.Register(c); err != nil {
			if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
				panic(err)
			}
		}
	}
}
