package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

type HTTPHandler struct {
	service *OrderService
}

func NewHTTPHandler(service *OrderService) *HTTPHandler {
	return &HTTPHandler{service: service}
}

func (h *HTTPHandler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /orders", h.handleCheckout)
	mux.HandleFunc("GET /orders/{id}", h.handleGetOrder)
	mux.HandleFunc("GET /orders", h.handleListOrders)
	return mux
}

func (h *HTTPHandler) handleCheckout(w http.ResponseWriter, r *http.Request) {
	var dto CreateOrderDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	order, err := h.service.Checkout(r.Context(), dto)
	if err != nil {
		if errors.Is(err, ErrInvalidOrder) || errors.Is(err, ErrInvalidItem) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, ErrPaymentFailed) {
			http.Error(w, err.Error(), http.StatusPaymentRequired)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(order)
}

func (h *HTTPHandler) handleGetOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	order, err := h.service.GetOrder(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrOrderNotFound) {
			http.Error(w, "order not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(order)
}

func (h *HTTPHandler) handleListOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.service.ListOrders(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(orders)
}
