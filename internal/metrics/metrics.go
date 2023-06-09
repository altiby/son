package metrics

import (
	metrics "github.com/altiby/son/internal/gometrics"
	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"strings"
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

func InstrumentChiRouter(router *chi.Mux, appName string) (*chi.Mux, error) {
	hm := metrics.NewHandlerMetrics(appName)

	rx := chi.NewRouter()
	rx.Use(router.Middlewares()...)
	err := chi.Walk(router, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		route = strings.TrimSuffix(strings.Replace(route, "/*/", "/", -1), "/")

		rx.With(middlewares...).Method(method, route, hm.Handler(method, route, handler))

		return nil
	})

	return rx, err
}
