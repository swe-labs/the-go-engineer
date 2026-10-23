package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config encapsulates validated application configuration parameters.
type Config struct {
	Port           int
	Environment    string // "development", "staging", "production"
	DatabaseURL    string
	JWTSecret      string
	WorkerCount    int
	LogLevel       string
	RequestTimeout time.Duration
}

// LoadConfig reads configuration from environment variables with fallback defaults and strict validation.
func LoadConfig() (*Config, error) {
	cfg := &Config{
		Port:           8080,
		Environment:    "development",
		DatabaseURL:    "postgres://opslane:opslane_secret@localhost:5432/opslane_dev?sslmode=disable",
		JWTSecret:      "opslane-dev-secret-key-32bytes!!",
		WorkerCount:    4,
		LogLevel:       "info",
		RequestTimeout: 30 * time.Second,
	}

	if p := os.Getenv("PORT"); p != "" {
		port, err := strconv.Atoi(p)
		if err != nil || port < 1024 || port > 65535 {
			return nil, fmt.Errorf("invalid PORT %q: must be integer between 1024 and 65535", p)
		}
		cfg.Port = port
	}

	if env := os.Getenv("ENV"); env != "" {
		validEnvs := map[string]bool{"development": true, "staging": true, "production": true}
		if !validEnvs[env] {
			return nil, fmt.Errorf("invalid ENV %q: must be development, staging, or production", env)
		}
		cfg.Environment = env
	}

	if db := os.Getenv("DB_URL"); db != "" {
		cfg.DatabaseURL = db
	}

	if jwt := os.Getenv("JWT_SECRET"); jwt != "" {
		if len(jwt) < 16 {
			return nil, errors.New("JWT_SECRET must be at least 16 characters for security")
		}
		cfg.JWTSecret = jwt
	}

	if wc := os.Getenv("WORKER_COUNT"); wc != "" {
		count, err := strconv.Atoi(wc)
		if err != nil || count < 1 || count > 64 {
			return nil, fmt.Errorf("invalid WORKER_COUNT %q: must be between 1 and 64", wc)
		}
		cfg.WorkerCount = count
	}

	if ll := os.Getenv("LOG_LEVEL"); ll != "" {
		cfg.LogLevel = strings.ToLower(ll)
	}

	if rt := os.Getenv("REQUEST_TIMEOUT"); rt != "" {
		dur, err := time.ParseDuration(rt)
		if err != nil {
			return nil, fmt.Errorf("invalid REQUEST_TIMEOUT %q: %w", rt, err)
		}
		cfg.RequestTimeout = dur
	}

	// In production, strictly enforce non-default secrets
	if cfg.Environment == "production" {
		if cfg.JWTSecret == "opslane-dev-secret-key-32bytes!!" {
			return nil, errors.New("production environment cannot use default insecure JWT_SECRET")
		}
	}

	return cfg, nil
}

// InitLogger constructs a structured slog logger based on the configured log level and format.
func InitLogger(levelStr string, isProd bool) *slog.Logger {
	var level slog.Level
	switch strings.ToLower(levelStr) {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{Level: level}
	var handler slog.Handler
	if isProd {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}

// Application manages the runtime lifecycle and dependencies of the Opslane service.
type Application struct {
	Config    *Config
	Logger    *slog.Logger
	StartTime time.Time
}

// NewApplication wires configuration and core dependencies.
func NewApplication(cfg *Config) *Application {
	isProd := cfg.Environment == "production"
	logger := InitLogger(cfg.LogLevel, isProd)

	return &Application{
		Config:    cfg,
		Logger:    logger,
		StartTime: time.Now().UTC(),
	}
}

// Routes sets up root health and system status routes.
func (app *Application) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", app.handleHealthCheck)
	mux.HandleFunc("GET /readyz", app.handleReadinessCheck)

	return app.recoverMiddleware(mux)
}

func (app *Application) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(app.StartTime).Round(time.Second).String()
	res := map[string]any{
		"status":      "ok",
		"service":     "opslane-backend",
		"environment": app.Config.Environment,
		"uptime":      uptime,
		"timestamp":   time.Now().UTC(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
}

func (app *Application) handleReadinessCheck(w http.ResponseWriter, r *http.Request) {
	res := map[string]any{
		"ready":     true,
		"workers":   app.Config.WorkerCount,
		"timestamp": time.Now().UTC(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
}

func (app *Application) recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				app.Logger.Error("server panic recovered", "error", rec)
				http.Error(w, `{"error":"internal_server_error"}`, http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "FATAL: Configuration validation failed: %v\n", err)
		os.Exit(1)
	}

	app := NewApplication(cfg)
	app.Logger.Info("Starting Opslane Flagship Backend Foundation",
		"port", cfg.Port,
		"env", cfg.Environment,
		"workers", cfg.WorkerCount,
	)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      app.Routes(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		app.Logger.Error("Server terminated unexpectedly", "error", err)
		os.Exit(1)
	}
}
