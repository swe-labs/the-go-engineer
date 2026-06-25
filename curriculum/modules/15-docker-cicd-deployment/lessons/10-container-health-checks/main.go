package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)

type HealthStatus struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Uptime    string `json:"uptime"`
}

var (
	healthy   atomic.Bool
	ready     atomic.Bool
	startTime time.Time
)

func init() {
	healthy.Store(false)
	ready.Store(false)
	startTime = time.Now()
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	if !healthy.Load() {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(HealthStatus{
			Status:    "unhealthy",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Uptime:    time.Since(startTime).String(),
		})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Uptime:    time.Since(startTime).String(),
	})
}

func readyzHandler(w http.ResponseWriter, r *http.Request) {
	if !ready.Load() {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{"status": "not ready"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
}

func main() {
	healthy.Store(true)

	go func() {
		time.Sleep(2 * time.Second)
		ready.Store(true)
		slog.Info("service is ready")
	}()

	http.HandleFunc("/healthz", healthzHandler)
	http.HandleFunc("/readyz", readyzHandler)

	addr := os.Getenv("LISTEN_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	slog.Info("server starting", "addr", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
