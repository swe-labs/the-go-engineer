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

	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/auth"
	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/config"
	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/db"
	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/events"
	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/handlers"
	paymentflow "github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/payment"
	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/services"
	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/workers"
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
	_ = orderPool.Start(context.Background())

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
	_ = paymentPool.Start(context.Background())

	isDraining := &atomic.Bool{}

	app := &handlers.Application{
		Logger:            logger,
		Store:             store,
		Orders:            orders,
		Payments:          payments,
		Tokens:            tokens,
		ServiceName:       cfg.App.Name,
		Environment:       cfg.App.Env,
		TrustedProxyCIDRs: cfg.HTTP.TrustedProxyCIDRs,
		IsDraining:        isDraining,
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

	idleConnsClosed := setupGracefulShutdown(server, logger, isDraining, bus, orderPool, paymentPool)

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("server stopped unexpectedly", slog.Any("error", err))
		os.Exit(1)
	}

	// Wait for graceful shutdown sequence to complete before returning and dropping
	// the database connection (via defer database.Close()).
	<-idleConnsClosed
	logger.Info("opslane server gracefully stopped")
}
