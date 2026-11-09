package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type probeState struct {
	mu         sync.RWMutex
	dbReady    bool
	cacheReady bool
}

func (ps *probeState) setDBReady(v bool) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.dbReady = v
}

func (ps *probeState) setCacheReady(v bool) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.cacheReady = v
}

func (ps *probeState) isHealthy() bool {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	return ps.dbReady && ps.cacheReady
}

func (ps *probeState) isReady() bool {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	return ps.dbReady
}

var state = &probeState{}

type healthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if !state.isHealthy() {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(healthResponse{
			Status:    "unhealthy",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(healthResponse{
		Status:    "ok",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

func readyzHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if !state.isReady() {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(healthResponse{
			Status:    "not ready",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(healthResponse{
		Status:    "ready",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

type ProbeCheck func() bool

func livenessHandler(check ProbeCheck) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if !check() {
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(healthResponse{Status: "unhealthy", Timestamp: time.Now().UTC().Format(time.RFC3339)})
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(healthResponse{Status: "ok", Timestamp: time.Now().UTC().Format(time.RFC3339)})
	}
}

func main() {
	state.setDBReady(true)
	state.setCacheReady(true)

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthzHandler)
	mux.HandleFunc("/readyz", readyzHandler)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

func ensureResponse(w http.ResponseWriter) {
	fmt.Fprintln(w, "ok")
}
