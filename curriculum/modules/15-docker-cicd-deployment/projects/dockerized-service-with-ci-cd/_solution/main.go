package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync/atomic"
	"syscall"
	"time"
)

// Injected during build time via -ldflags
var (
	Version   = "v1.0.0"
	CommitSHA = "dev"
)

type Config struct {
	Port         string
	Env          string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func LoadConfig() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	env := os.Getenv("ENV")
	if env == "" {
		env = "production"
	}
	return Config{
		Port:         port,
		Env:          env,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

type ServerState struct {
	isReady atomic.Bool
	started time.Time
}

func NewServerState() *ServerState {
	s := &ServerState{
		started: time.Now().UTC(),
	}
	s.isReady.Store(true)
	return s
}

func SetupRoutes(state *ServerState) http.Handler {
	mux := http.NewServeMux()

	// Liveness Probe: Kubelet / container orchestrator checks if process is alive
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	// Readiness Probe: Checks if server is prepared to route live customer traffic
	mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if !state.isReady.Load() {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte(`{"status":"not_ready"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ready"}`))
	})

	// Service Metadata & Diagnostic Endpoint
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

func runHealthCheck(addr string) int {
	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(fmt.Sprintf("http://127.0.0.1%s/healthz", addr))
	if err != nil || resp.StatusCode != http.StatusOK {
		return 1
	}
	return 0
}

func main() {
	checkHealth := flag.Bool("healthcheck", false, "run internal container healthcheck probe")
	flag.Parse()

	cfg := LoadConfig()

	if *checkHealth {
		os.Exit(runHealthCheck(":" + cfg.Port))
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	state := NewServerState()
	handler := SetupRoutes(state)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		slog.Info("server starting",
			slog.String("port", cfg.Port),
			slog.String("env", cfg.Env),
			slog.String("version", Version),
			slog.String("commit", CommitSHA),
		)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server encountered fatal error", "error", err)
			os.Exit(1)
		}
	}()

	<-shutdown
	slog.Info("shutdown signal received, commencing graceful termination")

	// Step 1: Mark not ready so orchestrator stops routing new traffic
	state.isReady.Store(false)

	// Step 2: Allow existing in-flight connections 5s to finish
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("graceful shutdown failed", "error", err)
	} else {
		slog.Info("clean server shutdown complete")
	}
}
