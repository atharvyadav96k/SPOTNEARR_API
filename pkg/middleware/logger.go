package middleware

import (
	"log"
	"net/http"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
	logger *log.Logger
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func RequestLogger(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK, logger: logger}
			next.ServeHTTP(rec, r)
			rec.logger.Printf("[%s] %s → %d (%s)", r.Method, r.URL.RequestURI(), rec.status, time.Since(start))
		})
	}
}
