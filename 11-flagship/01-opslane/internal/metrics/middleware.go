// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package metrics

import (
	"net/http"
	"time"
)

// statusRecorder (Struct): wraps http.ResponseWriter to capture the status code for metrics
type statusRecorder struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

// WriteHeader (Method): intercepts the status code for metric recording before passing to the client
func (sr *statusRecorder) WriteHeader(code int) {
	if !sr.wroteHeader {
		sr.status = code
		sr.wroteHeader = true
	}
	sr.ResponseWriter.WriteHeader(code)
}

// Write (Method): assumes default 200 OK status for metrics if WriteHeader was not explicitly called
func (sr *statusRecorder) Write(b []byte) (int, error) {
	if !sr.wroteHeader {
		sr.status = http.StatusOK
		sr.wroteHeader = true
	}
	return sr.ResponseWriter.Write(b)
}

// HTTPMetrics (Function): returns HTTP middleware recording request count, latency, and status class metrics
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
