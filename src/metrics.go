package main

import (
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP metrics
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// Database metrics
	dbOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "db_operations_total",
			Help: "Total number of database operations",
		},
		[]string{"operation", "table"},
	)

	dbOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_operation_duration_seconds",
			Help:    "Duration of database operations in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation", "table"},
	)

	// Business metrics
	eventsCreatedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "events_created_total",
			Help: "Total number of events created",
		},
	)

	eventsDeletedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "events_deleted_total",
			Help: "Total number of events deleted",
		},
	)

	activeEvents = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_events",
			Help: "Number of active events in the system",
		},
	)

	// System metrics
	goroutines = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "goroutines",
			Help: "Number of goroutines",
		},
	)

	memoryAlloc = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "memory_alloc_bytes",
			Help: "Allocated memory in bytes",
		},
	)
)

// MetricsMiddleware adds Prometheus metrics to HTTP handlers
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		path := r.URL.Path
		method := r.Method

		// Create a custom response writer to capture the status code
		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(rw, r)

		// Record metrics
		duration := time.Since(start).Seconds()
		httpRequestsTotal.WithLabelValues(method, path, strconv.Itoa(rw.statusCode)).Inc()
		httpRequestDuration.WithLabelValues(method, path).Observe(duration)
	})
}

// RecordDBOperation records database operation metrics
func RecordDBOperation(operation, table string, duration time.Duration) {
	dbOperationsTotal.WithLabelValues(operation, table).Inc()
	dbOperationDuration.WithLabelValues(operation, table).Observe(duration.Seconds())
}

// UpdateSystemMetrics updates system-level metrics
func UpdateSystemMetrics() {
	goroutines.Set(float64(runtime.NumGoroutine()))
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	memoryAlloc.Set(float64(m.Alloc))
}
