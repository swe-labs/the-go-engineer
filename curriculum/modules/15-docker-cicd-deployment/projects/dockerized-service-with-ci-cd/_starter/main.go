package main

import (
	"encoding/json"
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

var (
	Version   = "v0.1.0-dev"
	CommitSHA = "dirty"
)

type Config struct {
	Port string
	Env  string
}

func LoadConfig() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return Config{
		Port: port,
		Env:  "development",
	}
}

type ServerState struct {
	// TODO: Task 2 - Track readiness state with atomic.Bool for clean drain shutdown
	started time.Time
}

func NewServerState() *ServerState {
	return &ServerState{
		started: time.Now().UTC(),
	}
}

func SetupRoutes(state *ServerState) http.Handler {
	mux := http.NewServeMux()

	// Liveness Probe
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	// TODO: Task 2 - Implement /readyz readiness probe endpoint returning 200 or 503
	mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ready"}`))
	})

	// Metadata
	mux.HandleFunc("GET /api/v1/info", func(w http.ResponseWriter, r *http.Request) {
		info := map[string]any{
			"version":    Version,
			"commit_sha": CommitSHA,
			"go_version": runtime.Version(),
			"arch":       runtime.GOARCH,
			"os":         runtime.GOOS,
			"uptime_sec": int(time.Since(state.started).Seconds()),
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(info)
	})

	return mux
}

// TODO: Task 2 - Implement container internal healthcheck probe function
func runHealthCheck(addr string) int {
	return 0
}

func main() {
	checkHealth := flag.Bool("healthcheck", false, "run internal container healthcheck probe")
	flag.Parse()

	cfg := LoadConfig()

	if *checkHealth {
		os.Exit(runHealthCheck(":" + cfg.Port))
	}

	state := NewServerState()
	handler := SetupRoutes(state)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: handler,
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		slog.Info("starter server running on :" + cfg.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "error", err)
		}
	}()

	<-shutdown
	slog.Info("stopping server")
}
