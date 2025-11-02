package metrics

import (
	"log/slog"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	RequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "app",
		Name:      "requests_total",
		Help:      "Total amount of requests",
	}, []string{"method", "path", "status code"})
	RequestDuration = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: "app",
		Name:      "request_duration",
		Help:      "Summary of duration for certain path, method, status code",
		Objectives: map[float64]float64{
			0.5:  0.01,
			0.9:  0.005,
			0.99: 0.001,
		},
	}, []string{"method", "path", "status code"})
)

func SetupMetrics(host string) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	slog.Debug("Starting metrics...", "host", host)
	return http.ListenAndServe(host, mux)
}
