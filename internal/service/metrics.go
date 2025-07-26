package service

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"sweng-task/internal/config"
)

// This part is for fiber-go analytics collect
var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"path", "method", "code"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Latency of HTTP requests in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path", "method"},
	)

	httpRequestsInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Current number of in-flight HTTP requests.",
		},
	)

	httpRequestSizeBytes = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_size_bytes",
			Help:    "Size of HTTP requests in bytes.",
			Buckets: prometheus.ExponentialBuckets(100, 10, 6),
		},
		[]string{"path", "method"},
	)

	httpResponseSizeBytes = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "Size of HTTP responses in bytes.",
			Buckets: prometheus.ExponentialBuckets(100, 10, 6),
		},
		[]string{"path", "method"},
	)
)

type Metrics struct {
	logs *zap.SugaredLogger
	cfg  *config.Config
}

func NewMetricsService(log *zap.SugaredLogger, cf *config.Config) *Metrics {
	return &Metrics{
		logs: log,
		cfg:  cf,
	}
}

// NewPrometheusMiddleware creates a new Fiber middleware for Prometheus metrics.
func (m *Metrics) NewPrometheusMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		httpRequestsInFlight.Inc()
		defer httpRequestsInFlight.Dec()

		reqSize := float64(c.Request().Header.ContentLength())
		if reqSize > 0 {
			httpRequestSizeBytes.WithLabelValues(c.Path(), c.Method()).Observe(reqSize)
		}

		timer := prometheus.NewTimer(httpRequestDuration.WithLabelValues(c.Path(), c.Method()))
		err := c.Next()

		// After the handler has finished, record the duration
		timer.ObserveDuration()

		respSize := float64(len(c.Response().Body()))
		httpResponseSizeBytes.WithLabelValues(c.Path(), c.Method()).Observe(respSize)

		statusCode := strconv.Itoa(c.Response().StatusCode())
		httpRequestsTotal.WithLabelValues(c.Path(), c.Method(), statusCode).Inc()

		return err
	}
}

func (m *Metrics) Start() {
	metricsServer := http.NewServeMux()
	metricsServer.Handle("/metrics", promhttp.Handler())

	m.logs.Infof("Starting metrics server on :%d", m.cfg.Metrics.Port)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", m.cfg.Metrics.Port), metricsServer); err != nil {
		m.logs.Fatalf("Failed to start metrics server: %v", err)
	}
}
