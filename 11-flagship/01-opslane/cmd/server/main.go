// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/swe-labs/the-go-engineer/11-flagship/01-opslane/internal/auth"
	"github.com/swe-labs/the-go-engineer/11-flagship/01-opslane/internal/config"
	"github.com/swe-labs/the-go-engineer/11-flagship/01-opslane/internal/db"
	"github.com/swe-labs/the-go-engineer/11-flagship/01-opslane/internal/events"
	"github.com/swe-labs/the-go-engineer/11-flagship/01-opslane/internal/handlers"
	"github.com/swe-labs/the-go-engineer/11-flagship/01-opslane/internal/metrics"
	"github.com/swe-labs/the-go-engineer/11-flagship/01-opslane/internal/otel"
	paymentflow "github.com/swe-labs/the-go-engineer/11-flagship/01-opslane/internal/payment"
	"github.com/swe-labs/the-go-engineer/11-flagship/01-opslane/internal/ratelimit"
	"github.com/swe-labs/the-go-engineer/11-flagship/01-opslane/internal/services"
	"github.com/swe-labs/the-go-engineer/11-flagship/01-opslane/internal/workers"
)

const startupDatabaseTimeout = 10 * time.Second

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.New(slog.NewTextHandler(os.Stderr, nil)).Error("failed to load config", slog.Any("error", err))
		os.Exit(1)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.App.LogLevel,
	}))

	startupCtx, cancelStartup := context.WithTimeout(context.Background(), startupDatabaseTimeout)
	defer cancelStartup()

	database, err := db.Open(startupCtx, cfg.Database)
	if err != nil {
		logger.Error("failed to open database", slog.Any("error", err))
		os.Exit(1)
	}
	defer database.Close()

	if err := db.Migrate(startupCtx, database); err != nil {
		logger.Error("failed to apply database migrations", slog.Any("error", err))
		os.Exit(1)
	}

	var otelCfg otel.Config
	otelCfg.FromEnv()
	tracer := otel.New(otelCfg, logger)
	defer tracer.Stop()

	rateLimiter := ratelimit.New(ratelimit.Config{
		RequestsPerSecond: 2,
		BurstSize:         120,
		WindowSeconds:     60,
		DB:                database,
	})

	tokens, err := auth.NewTokenManager(cfg.Auth.TokenSecret, cfg.Auth.TokenIssuer, cfg.Auth.TokenTTL)
	if err != nil {
		logger.Error("failed to initialize auth tokens", slog.Any("error", err))
		os.Exit(1)
	}

	store := db.NewStore(database)
	orders := services.NewOrderService(store, services.NewNoopInventoryCoordinator())
	payments := services.NewPaymentService(store, store, orders, paymentflow.NewSimulatedGateway(), services.PaymentServiceOptions{})

	// Initialize background systems
	bus := events.NewBus(1000)

	// Create root application context for workers
	ctx, cancelApp := context.WithCancel(context.Background())
	defer cancelApp()

	orderPool, err := workers.NewPool(workers.PoolConfig{
		Name:      "orders",
		Workers:   3,
		QueueSize: 500,
		Handler:   workers.OrderProcessor{Workflow: orders}.Handle,
	})
	if err != nil {
		logger.Error("failed to create order worker pool", slog.Any("error", err))
		os.Exit(1)
	}
	if err := orderPool.Start(ctx); err != nil {
		logger.Error("failed to start order worker pool", slog.Any("error", err))
		os.Exit(1)
	}

	paymentPool, err := workers.NewPool(workers.PoolConfig{
		Name:      "payments",
		Workers:   3,
		QueueSize: 500,
		Handler:   workers.PaymentProcessor{Workflow: payments}.Handle,
	})
	if err != nil {
		logger.Error("failed to create payment worker pool", slog.Any("error", err))
		os.Exit(1)
	}
	if err := paymentPool.Start(ctx); err != nil {
		logger.Error("failed to start payment worker pool", slog.Any("error", err))
		os.Exit(1)
	}

	isDraining := &atomic.Bool{}
	appMetrics := metrics.NewAppMetrics()

	app := &handlers.Application{
		Logger:            logger,
		Metrics:           appMetrics,
		Store:             store,
		Orders:            orders,
		Payments:          payments,
		Tokens:            tokens,
		ServiceName:       cfg.App.Name,
		Environment:       cfg.App.Env,
		TrustedProxyCIDRs: cfg.HTTP.TrustedProxyCIDRs,
		IsDraining:        isDraining,
		RateLimiter:       rateLimiter,
		Tracer:            tracer,
	}

	server := &http.Server{
		Addr:              cfg.HTTP.Address,
		Handler:           app.Routes(),
		ReadHeaderTimeout: cfg.HTTP.ReadHeaderTimeout,
		ReadTimeout:       cfg.HTTP.ReadTimeout,
		WriteTimeout:      cfg.HTTP.WriteTimeout,
		IdleTimeout:       cfg.HTTP.IdleTimeout,
	}

	logger.Info("starting opslane server",
		slog.String("env", cfg.App.Env),
		slog.String("addr", cfg.HTTP.Address),
		slog.String("database", "postgresql"),
	)

	idleConnsClosed := setupGracefulShutdown(server, cfg.HTTP.ShutdownTimeout, logger, isDraining, bus, cancelApp, orderPool, paymentPool)

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("server stopped unexpectedly", slog.Any("error", err))
		os.Exit(1)
	}

	// Wait for graceful shutdown sequence to complete before returning and dropping
	// the database connection (via defer database.Close()).
	<-idleConnsClosed
	logger.Info("opslane server gracefully stopped")
}
