// Collecting internal and runtime application metrics.
package metrics

import (
	"runtime"
	"time"
)

// ExportRunningGoRoutines – start process to record periodically number of goroutines.
func ExportRunningGoRoutines(m Metrics) {
	t := time.NewTicker(time.Second)
	go exportGoroutines(m, t)
}

func exportGoroutines(m Metrics, ticker *time.Ticker) {
	for range ticker.C {
		m.GaugeSet("goroutines.count", []string{}, float64(runtime.NumGoroutine()))
	}
}
