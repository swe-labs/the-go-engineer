// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package metrics

import (
	"fmt"
	"io"
	"net/http"
)

// PrometheusHandler returns an HTTP handler that exposes metrics
// in Prometheus text format. This endpoint should be mounted at
// /metrics for Prometheus scraper compatibility.
func PrometheusHandler(m *AppMetrics) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")

		rep := &prometheusWriter{Writer: w}

		fmt.Fprintf(rep, "# HELP opslane_http_requests_total Total HTTP requests\n")
		fmt.Fprintf(rep, "# TYPE opslane_http_requests_total counter\n")
		fmt.Fprintf(rep, "opslane_http_requests_total %d\n\n", m.HTTPRequestsTotal.Value())

		fmt.Fprintf(rep, "# HELP opslane_http_request_duration_seconds HTTP request duration in seconds\n")
		fmt.Fprintf(rep, "# TYPE opslane_http_request_duration_seconds histogram\n")
		printHistogramProm(rep, "opslane_http_request_duration_seconds", m.HTTPRequestDuration)

		fmt.Fprintf(rep, "# HELP opslane_http_responses_total HTTP responses by status class\n")
		fmt.Fprintf(rep, "# TYPE opslane_http_responses_total counter\n")
		fmt.Fprintf(rep, "opslane_http_responses_total{status=\"2xx\"} %d\n", m.HTTPResponses2xx.Value())
		fmt.Fprintf(rep, "opslane_http_responses_total{status=\"4xx\"} %d\n", m.HTTPResponses4xx.Value())
		fmt.Fprintf(rep, "opslane_http_responses_total{status=\"5xx\"} %d\n\n", m.HTTPResponses5xx.Value())

		fmt.Fprintf(rep, "# HELP opslane_cache_hits_total Cache hit count\n")
		fmt.Fprintf(rep, "# TYPE opslane_cache_hits_total counter\n")
		fmt.Fprintf(rep, "opslane_cache_hits_total %d\n\n", m.CacheHits.Value())

		fmt.Fprintf(rep, "# HELP opslane_cache_misses_total Cache miss count\n")
		fmt.Fprintf(rep, "# TYPE opslane_cache_misses_total counter\n")
		fmt.Fprintf(rep, "opslane_cache_misses_total %d\n\n", m.CacheMisses.Value())

		fmt.Fprintf(rep, "# HELP opslane_worker_jobs_total Worker job count\n")
		fmt.Fprintf(rep, "# TYPE opslane_worker_jobs_total counter\n")
		fmt.Fprintf(rep, "opslane_worker_jobs_total{status=\"processed\"} %d\n", m.WorkerJobsProcessed.Value())
		fmt.Fprintf(rep, "opslane_worker_jobs_total{status=\"failed\"} %d\n\n", m.WorkerJobsFailed.Value())
	}
}

type prometheusWriter struct {
	io.Writer
}

func printHistogramProm(w io.Writer, name string, h *Histogram) {
	snap := h.Snapshot()

	fmt.Fprintf(w, "%s_count %d\n", name, int64(snap.Count))
	fmt.Fprintf(w, "%s_sum %f\n", name, snap.Sum)

	le := float64(0)
	fmt.Fprintf(w, "%s_bucket{le=\"%g\"} 0\n", name, le)
	for i, bound := range snap.Bounds {
		le = bound
		cumulative := int64(0)
		for j := 0; j <= i; j++ {
			cumulative += snap.Buckets[j]
		}
		fmt.Fprintf(w, "%s_bucket{le=\"%g\"} %d\n", name, bound, cumulative)
	}
	// +Inf bucket = total count
	fmt.Fprintf(w, "%s_bucket{le=\"+Inf\"} %d\n\n", name, snap.Count)
}