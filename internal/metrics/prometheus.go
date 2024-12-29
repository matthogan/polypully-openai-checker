package metrics

import (
	"github.com/codejago/polypully-openai-checker/internal/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log/slog"
	"net/http"
)

var (
	// global configuration
	metrix *config.Metrics
	// Metrics definitions
	startupCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "service_startup_total",
			Help: "Total number of times the service has started up",
		},
	)

	shutdownCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "service_shutdown_total",
			Help: "Total number of times the service has shut down.",
		},
	)

	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests received, labeled by method and status.",
		},
		[]string{"method", "status"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of response time for handler in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)
)

func StartMetrics(config *config.Metrics) {
	metrix = config
	if metrix.Enabled {
		slog.Info("metrics enabled ✔")
		// Register metrics with the Prometheus default registry
		prometheus.MustRegister(startupCounter)
		prometheus.MustRegister(shutdownCounter)
		prometheus.MustRegister(requestCounter)
		prometheus.MustRegister(requestDuration)
		// Serve Prometheus metrics on /metrics
		http.Handle(metrix.ContextRoot, promhttp.Handler())
		// Start the metrics server
		go func() {
			slog.Info("starting metrics server on", "address", metrix.Localhost)
			if err := http.ListenAndServe(metrix.Localhost, nil); err != nil {
				slog.Error("error starting metrics server", "error", err)
			}
		}()
		startupCounter.Inc()
	} else {
		slog.Info("metrics enabled 〤")
	}
}

func StopMetrics() {
	if metrix.Enabled {
		shutdownCounter.Inc()
	}
}

func IncRequestCounter(method, status string) {
	if metrix.Enabled {
		requestCounter.WithLabelValues(method, status).Inc()
	}
}

func SetRequestDuration(method, status string, duration float64) {
	if metrix.Enabled {
		requestDuration.WithLabelValues(method, status).Observe(duration)
	}
}
