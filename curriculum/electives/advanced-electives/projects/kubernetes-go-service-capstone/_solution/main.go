// Command orders-api is the reference implementation for the Kubernetes Go
// Service Capstone. It is a small, self-contained HTTP service that
// demonstrates every object the capstone requires: a version endpoint that
// proves which revision answered, liveness and readiness probes, a bounded
// graceful drain, configuration loaded from the environment and a mounted
// Secret file, and a memory budget derived from the container limit.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

// currentSchema identifies the /version response shape. Bump it only when a
// field is removed or changes meaning; adding a field keeps the schema version
// so an older client keeps decoding a newer response.
const currentSchema = "version.v1"

// VersionInfo is the payload returned by /version.
type VersionInfo struct {
	Schema      string `json:"schema"`
	Version     string `json:"version"`
	Commit      string `json:"commit"`
	BuildTime   string `json:"build_time"`
	Modified    bool   `json:"modified"`
	Environment string `json:"environment"`
}

// Config is the validated runtime configuration.
type Config struct {
	Addr        string
	Environment string
	LogLevel    string
	MaxConns    int
	DatabaseURL string
	MemLimit    int64
}

// ConfigError reports the first invalid configuration value. Returning an error
// rather than exiting inside the loader keeps the loader testable.
type ConfigError struct{ Field, Reason string }

func (e *ConfigError) Error() string {
	return fmt.Sprintf("config error: %s %s", e.Field, e.Reason)
}

// LoadConfig reads configuration from the environment and from the file named
// by APP_DATABASE_URL_FILE. getenv and readFile are injected so the loader can
// be tested without touching the real process environment or filesystem.
func LoadConfig(getenv func(string) string, readFile func(string) ([]byte, error)) (Config, error) {
	c := Config{
		Addr:        ":8080",
		Environment: getenv("APP_ENV"),
		LogLevel:    getenv("APP_LOG_LEVEL"),
		MemLimit:    -1,
	}
	if c.Environment == "" {
		c.Environment = "local"
	}
	if c.LogLevel == "" {
		c.LogLevel = "info"
	}
	switch c.LogLevel {
	case "debug", "info", "warn", "error":
	default:
		return Config{}, &ConfigError{Field: "APP_LOG_LEVEL", Reason: "must be one of debug, info, warn, error"}
	}
	if v := getenv("APP_ADDR"); v != "" {
		c.Addr = v
	}

	maxConns := getenv("APP_MAX_CONNS")
	if strings.TrimSpace(maxConns) == "" {
		return Config{}, &ConfigError{Field: "APP_MAX_CONNS", Reason: "is required"}
	}
	n, err := strconv.Atoi(maxConns)
	if err != nil || n <= 0 {
		return Config{}, &ConfigError{Field: "APP_MAX_CONNS", Reason: "must be a positive integer"}
	}
	c.MaxConns = n

	path := getenv("APP_DATABASE_URL_FILE")
	if strings.TrimSpace(path) == "" {
		return Config{}, &ConfigError{Field: "APP_DATABASE_URL_FILE", Reason: "is required"}
	}
	raw, err := readFile(path)
	if err != nil {
		return Config{}, &ConfigError{Field: "APP_DATABASE_URL_FILE", Reason: fmt.Sprintf("cannot read %s: %v", path, err)}
	}
	c.DatabaseURL = strings.TrimSpace(string(raw))
	if c.DatabaseURL == "" {
		return Config{}, &ConfigError{Field: "APP_DATABASE_URL_FILE", Reason: "contains an empty value"}
	}

	if v := getenv("APP_MEMORY_LIMIT_BYTES"); v != "" {
		limit, err := strconv.ParseInt(v, 10, 64)
		if err != nil || limit <= 0 {
			return Config{}, &ConfigError{Field: "APP_MEMORY_LIMIT_BYTES", Reason: "must be a positive integer"}
		}
		c.MemLimit = limit
	}
	return c, nil
}

// ComputeBudget reserves headroom below the container memory limit so the
// garbage collector reacts before the kernel OOM killer. It returns the soft
// heap limit to hand to debug.SetMemoryLimit, or -1 when no limit is known.
func ComputeBudget(containerLimit int64) int64 {
	if containerLimit <= 0 {
		return -1
	}
	return containerLimit - containerLimit/8
}

// Server owns the readiness state and the HTTP listener.
type Server struct {
	ready     atomic.Bool
	inflight  sync.WaitGroup
	env       string
	buildInfo *debug.BuildInfo
	handler   http.Handler
}

// NewServer returns a server in the not-ready state. Call Ready to open traffic.
func NewServer(env string, bi *debug.BuildInfo) *Server {
	s := &Server{env: env, buildInfo: bi}
	s.handler = s.routes()
	return s
}

// Ready marks the server ready to receive traffic.
func (s *Server) Ready() { s.ready.Store(true) }

// NotReady marks the server as draining so /readyz returns 503. It must be
// called before the listener is closed so the endpoint controller removes this
// Pod from the Service endpoints.
func (s *Server) NotReady() { s.ready.Store(false) }

func (s *Server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok\n"))
	})
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		if !s.ready.Load() {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte("not ready\n"))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ready\n"))
	})
	mux.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(BuildVersion(s.env, s.buildInfo))
	})
	mux.HandleFunc("/work", func(w http.ResponseWriter, r *http.Request) {
		// A tracked handler: the in-flight counter lets Shutdown wait for the
		// requests the server already accepted.
		s.inflight.Add(1)
		defer s.inflight.Done()
		select {
		case <-time.After(100 * time.Millisecond):
		case <-r.Context().Done():
		}
		w.WriteHeader(http.StatusOK)
	})
	return mux
}

// Handler exposes the routes for tests.
func (s *Server) Handler() http.Handler { return s.handler }

// BuildVersion assembles a VersionInfo from the build metadata the Go toolchain
// embeds and from the runtime environment. buildInfo may be nil.
func BuildVersion(env string, bi *debug.BuildInfo) VersionInfo {
	v := VersionInfo{
		Schema:      currentSchema,
		Version:     "(devel)",
		Commit:      "unknown",
		BuildTime:   "unknown",
		Environment: env,
	}
	if v.Environment == "" {
		v.Environment = "local"
	}
	if bi != nil {
		if bi.Main.Version != "" {
			v.Version = bi.Main.Version
		}
		for _, s := range bi.Settings {
			switch s.Key {
			case "vcs.revision":
				if s.Value != "" {
					v.Commit = s.Value
				}
			case "vcs.time":
				if s.Value != "" {
					v.BuildTime = s.Value
				}
			case "vcs.modified":
				v.Modified = s.Value == "true"
			}
		}
	}
	return v
}

// Run serves traffic on ln and drains on the provided context. It returns after
// the listener is closed and every in-flight request has finished, or when the
// drain window expires.
func (s *Server) Run(ctx context.Context, ln net.Listener, drain time.Duration) error {
	srv := &http.Server{Handler: s.handler}

	s.ready.Store(true)

	errCh := make(chan error, 1)
	go func() { errCh <- srv.Serve(ln) }()

	select {
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	case <-ctx.Done():
	}

	// Drain: stop routing new traffic, then close the listener and wait for the
	// requests already accepted to finish within the window.
	s.NotReady()
	shCtx, cancel := context.WithTimeout(context.Background(), drain)
	defer cancel()
	if err := srv.Shutdown(shCtx); err != nil {
		_ = srv.Close()
		return fmt.Errorf("drain exceeded %s: %w", drain, err)
	}
	done := make(chan struct{})
	go func() { s.inflight.Wait(); close(done) }()
	select {
	case <-done:
	case <-shCtx.Done():
		return fmt.Errorf("in-flight requests did not finish within %s", drain)
	}
	return nil
}

func main() {
	cfg, err := LoadConfig(os.Getenv, os.ReadFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if soft := ComputeBudget(cfg.MemLimit); soft > 0 {
		debug.SetMemoryLimit(soft)
	}

	bi, ok := debug.ReadBuildInfo()
	if !ok {
		bi = nil
	}

	fmt.Printf("listening on %s version=%s env=%s gomaxprocs=%d max_conns=%d\n",
		cfg.Addr, BuildVersion(cfg.Environment, bi).Version, cfg.Environment, runtime.GOMAXPROCS(0), cfg.MaxConns)

	s := NewServer(cfg.Environment, bi)

	ln, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "listen error: %v\n", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer stop()

	if err := s.Run(ctx, ln, 10*time.Second); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}
