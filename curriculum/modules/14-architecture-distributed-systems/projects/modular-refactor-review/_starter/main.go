package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

// ============================================================================
// LEGACY MONOLITH (BEFORE REFACTORING)
// ============================================================================
// Anti-patterns present in this starter:
// 1. Global mutable state (DB map guarded by global lock, not injectable).
// 2. Concrete hardcoded dependencies (cannot mock payment or notification in tests).
// 3. Mixed concerns: HTTP handler parses JSON, performs business calculation,
//    applies discount, mutates database, and sends notification.
// 4. Untestable via isolated unit tests without mutating global variables.
//
// TODO: Refactor this monolith into Hexagonal / Layered Architecture:
//   Task 1: Define pure domain models (Order, OrderItem, CalculateTotals) in domain.go
//   Task 2: Define port interfaces (OrderRepository, PaymentGateway, NotificationService) in ports.go
//   Task 3: Implement domain service with dependency injection in service.go
//   Task 4: Implement adapters (in-memory repo, mock payment, mock notifier, HTTP handler)
// ============================================================================

var (
	globalMu sync.RWMutex
	globalDB = make(map[string]map[string]any)
)

func generateLegacyID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func legacyCheckoutHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Customer string `json:"customer"`
		Email    string `json:"email"`
		Items    []struct {
			SKU        string `json:"sku"`
			Quantity   int    `json:"quantity"`
			PriceCents int    `json:"price_cents"`
		} `json:"items"`
		PaymentToken string `json:"payment_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(body.Customer) == "" || strings.TrimSpace(body.Email) == "" || len(body.Items) == 0 {
		http.Error(w, "invalid order data", http.StatusBadRequest)
		return
	}

	subtotal := 0
	for _, it := range body.Items {
		if strings.TrimSpace(it.SKU) == "" || it.Quantity <= 0 || it.PriceCents <= 0 {
			http.Error(w, "invalid item data", http.StatusBadRequest)
			return
		}
		subtotal += it.PriceCents * it.Quantity
	}

	discount := 0
	if subtotal >= 10000 {
		discount = subtotal / 10
	}
	total := subtotal - discount

	// Hardcoded payment check
	if body.PaymentToken == "fail_token" || body.PaymentToken == "" {
		http.Error(w, "payment failed", http.StatusPaymentRequired)
		return
	}
	paymentRef := fmt.Sprintf("ch_%s", generateLegacyID()[:8])

	orderID := generateLegacyID()
	record := map[string]any{
		"id":             orderID,
		"customer":       body.Customer,
		"email":          body.Email,
		"subtotal_cents": subtotal,
		"discount_cents": discount,
		"total_cents":    total,
		"status":         "PAID",
		"payment_ref":    paymentRef,
		"created_at":     time.Now().UTC().Format(time.RFC3339),
	}

	globalMu.Lock()
	globalDB[orderID] = record
	globalMu.Unlock()

	// Direct notification logging
	slog.Info("sent confirmation email", "to", body.Email, "order_id", orderID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(record)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /orders", legacyCheckoutHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		slog.Info("legacy monolith running on :8080")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "error", err)
		}
	}()

	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = server.Shutdown(ctx)
}
