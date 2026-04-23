// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package main

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/config"
)

type healthResponse struct {
	Service string `json:"service"`
	Env     string `json:"env"`
	Status  string `json:"status"`
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.New(slog.NewTextHandler(os.Stderr, nil)).Error("failed to load config", slog.Any("error", err))
		os.Exit(1)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.App.LogLevel,
	}))

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{
			"service": cfg.App.Name,
			"env":     cfg.App.Env,
			"message": "Opslane foundation shell is running",
		})
	})
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, healthResponse{
			Service: cfg.App.Name,
			Env:     cfg.App.Env,
			Status:  "ok",
		})
	})

	server := &http.Server{
		Addr:              cfg.HTTP.Address,
		Handler:           mux,
		ReadHeaderTimeout: cfg.HTTP.ReadHeaderTimeout,
		ReadTimeout:       cfg.HTTP.ReadTimeout,
		WriteTimeout:      cfg.HTTP.WriteTimeout,
		IdleTimeout:       cfg.HTTP.IdleTimeout,
	}

	logger.Info("starting opslane server",
		slog.String("env", cfg.App.Env),
		slog.String("addr", cfg.HTTP.Address),
	)

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("server stopped unexpectedly", slog.Any("error", err))
		os.Exit(1)
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
