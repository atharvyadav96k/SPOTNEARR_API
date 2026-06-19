package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "spotnearr_http_requests_total",
		Help: "Total HTTP requests by service, method, path, and status code.",
	}, []string{"service", "method", "path", "status"})

	httpRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "spotnearr_http_request_duration_seconds",
		Help:    "HTTP request latency in seconds.",
		Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
	}, []string{"service", "method", "path"})

	httpRequestsInFlight = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "spotnearr_http_requests_in_flight",
		Help: "Number of currently in-flight HTTP requests.",
	}, []string{"service"})
)

// PrometheusMetrics returns a middleware that records HTTP metrics per service.
// Uses the matched Gorilla Mux route template as the path label to avoid cardinality explosion.
func PrometheusMetrics(service string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/metrics" {
				next.ServeHTTP(w, r)
				return
			}

			path := routeTemplate(r)

			httpRequestsInFlight.WithLabelValues(service).Inc()
			defer httpRequestsInFlight.WithLabelValues(service).Dec()

			rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
			start := time.Now()
			next.ServeHTTP(rec, r)

			httpRequestsTotal.WithLabelValues(service, r.Method, path, strconv.Itoa(rec.status)).Inc()
			httpRequestDuration.WithLabelValues(service, r.Method, path).Observe(time.Since(start).Seconds())
		})
	}
}

// MetricsHandler returns the Prometheus HTTP handler for the /metrics endpoint.
func MetricsHandler() http.Handler {
	return promhttp.Handler()
}

// routeTemplate extracts the matched Gorilla Mux route template (e.g. /api/v1/users/{userId}/profile).
// Falls back to the raw request path if no route is matched yet.
func routeTemplate(r *http.Request) string {
	if route := mux.CurrentRoute(r); route != nil {
		if tmpl, err := route.GetPathTemplate(); err == nil {
			return tmpl
		}
	}
	return r.URL.Path
}
