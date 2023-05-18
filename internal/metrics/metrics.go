package metrics

import (
	metrics "github.com/altiby/son/internal/gometrics"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	defaultAddr = ":39901"
	defaultPath = "/metrics"
)

// NewPrometheus is a workaround function that creates additional server
// to host metrics, metrics library uses http.Handle, which does not export anything
// in case we will not use DefaultServerMux.
func NewPrometheus(appName string) (metrics.Metrics, error) {
	mux := http.NewServeMux()
	mux.Handle(defaultPath, promhttp.Handler())

	server := http.Server{
		Addr:    defaultAddr,
		Handler: mux,
	}

	errCh := make(chan error)

	go func() {
		errCh <- server.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		return nil, err
	default:
		return metrics.NewPrometheusMetrics(appName), nil
	}
}
