package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// API Standard Envelopes

type APIResponse struct {
	Data any `json:"data,omitempty"`
}

type FieldError struct {
	Field string `json:"field"`
	Issue string `json:"issue"`
}

type APIErrorEnvelope struct {
	Error struct {
		Code    string       `json:"code"`
		Message string       `json:"message"`
		Details []FieldError `json:"details,omitempty"`
	} `json:"error"`
}

// OrderDTO represents an incoming or outgoing order payload.
type OrderDTO struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Amount    float64   `json:"amount"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateOrderRequest struct {
	Title    string  `json:"title"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

func (req *CreateOrderRequest) Validate() []FieldError {
	var errs []FieldError
	if strings.TrimSpace(req.Title) == "" {
		errs = append(errs, FieldError{Field: "title", Issue: "title cannot be blank"})
	} else if len(req.Title) > 100 {
		errs = append(errs, FieldError{Field: "title", Issue: "title cannot exceed 100 characters"})
	}

	if req.Amount <= 0 {
		errs = append(errs, FieldError{Field: "amount", Issue: "amount must be greater than zero"})
	}

	curr := strings.ToUpper(strings.TrimSpace(req.Currency))
	if curr == "" {
		req.Currency = "USD"
	} else if curr != "USD" && curr != "EUR" && curr != "GBP" {
		errs = append(errs, FieldError{Field: "currency", Issue: "currency must be one of: USD, EUR, GBP"})
	}

	return errs
}

// APIServer encapsulates in-memory store and routing logic.
type APIServer struct {
	mu     sync.RWMutex
	orders map[string]*OrderDTO
}

func NewAPIServer() *APIServer {
	return &APIServer{
		orders: make(map[string]*OrderDTO),
	}
}

// Router builds standard library HTTP routes with middleware.
func (s *APIServer) Router() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/orders", s.handleCreateOrder)
	mux.HandleFunc("GET /api/v1/orders/{id}", s.handleGetOrder)
	mux.HandleFunc("GET /api/v1/orders", s.handleListOrders)
	mux.HandleFunc("DELETE /api/v1/orders/{id}", s.handleDeleteOrder)

	return s.requestIDMiddleware(mux)
}

func (s *APIServer) requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			b := make([]byte, 8)
			_, _ = rand.Read(b)
			reqID = hex.EncodeToString(b)
		}
		w.Header().Set("X-Request-ID", reqID)
		next.ServeHTTP(w, r)
	})
}

func (s *APIServer) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	// Enforce 1MB payload ceiling
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	var req CreateOrderRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			s.respondError(w, http.StatusRequestEntityTooLarge, "PAYLOAD_TOO_LARGE", "Request body exceeds maximum size of 1MB", nil)
			return
		}
		s.respondError(w, http.StatusBadRequest, "MALFORMED_JSON", "Invalid or malformed JSON payload", nil)
		return
	}

	if fieldErrors := req.Validate(); len(fieldErrors) > 0 {
		s.respondError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Request validation failed", fieldErrors)
		return
	}

	s.mu.Lock()
	id := fmt.Sprintf("ord-%d", len(s.orders)+1)
	order := &OrderDTO{
		ID:        id,
		Title:     strings.TrimSpace(req.Title),
		Amount:    req.Amount,
		Currency:  req.Currency,
		CreatedAt: time.Now().UTC(),
	}
	s.orders[id] = order
	s.mu.Unlock()

	s.respondJSON(w, http.StatusCreated, order)
}

func (s *APIServer) handleGetOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	s.mu.RLock()
	order, exists := s.orders[id]
	s.mu.RUnlock()

	if !exists {
		s.respondError(w, http.StatusNotFound, "ORDER_NOT_FOUND", fmt.Sprintf("Order %q not found", id), nil)
		return
	}

	s.respondJSON(w, http.StatusOK, order)
}

func (s *APIServer) handleListOrders(w http.ResponseWriter, r *http.Request) {
	limit := 10
	if l := r.URL.Query().Get("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 && val <= 100 {
			limit = val
		}
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*OrderDTO
	count := 0
	for _, o := range s.orders {
		if count >= limit {
			break
		}
		result = append(result, o)
		count++
	}

	s.respondJSON(w, http.StatusOK, result)
}

func (s *APIServer) handleDeleteOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	s.mu.Lock()
	_, exists := s.orders[id]
	if exists {
		delete(s.orders, id)
	}
	s.mu.Unlock()

	if !exists {
		s.respondError(w, http.StatusNotFound, "ORDER_NOT_FOUND", fmt.Sprintf("Order %q not found", id), nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *APIServer) respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(APIResponse{Data: data})
}

func (s *APIServer) respondError(w http.ResponseWriter, status int, code, message string, details []FieldError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	env := APIErrorEnvelope{}
	env.Error.Code = code
	env.Error.Message = message
	env.Error.Details = details

	_ = json.NewEncoder(w).Encode(env)
}

func main() {
	server := NewAPIServer()
	fmt.Println("Opslane HTTP REST API listening on :8080...")
	_ = http.ListenAndServe(":8080", server.Router())
}
