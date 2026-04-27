// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package metrics

import (
	"net/http"
	"time"
)

// statusRecorder wraps http.ResponseWriter to capture the status code.
type statusRecorder struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (sr *statusRecorder) WriteHeader(code int) {
	if !sr.wroteHeader {
		sr.status = code
		sr.wroteHeader = true
	}
	sr.ResponseWriter.WriteHeader(code)
}

func (sr *statusRecorder) Write(b []byte) (int, error) {
	if !sr.wroteHeader {
		sr.status = http.StatusOK
		sr.wroteHeader = true
	}
	return sr.ResponseWriter.Write(b)
}

// HTTPMetrics returns middleware that records per-request metrics:
//
//   - http_requests_total: incremented for every request
//   - http_request_duration_seconds: latency histogram
//   - http_responses_2xx / 4xx / 5xx: status class counters
//
// This middleware should be placed early in the middleware chain so
// it wraps the full request lifecycle.
func HTTPMetrics(m *AppMetrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			recorder := &statusRecorder{
				ResponseWriter: w,
				status:         http.StatusOK,
			}

			next.ServeHTTP(recorder, r)

			m.HTTPRequestsTotal.Inc()
			m.HTTPRequestDuration.Observe(time.Since(start).Seconds())
			m.ClassifyStatus(recorder.status)
		})
	}
}
