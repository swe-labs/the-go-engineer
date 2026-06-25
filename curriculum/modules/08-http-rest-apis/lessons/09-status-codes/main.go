package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type Order struct {
	ID     string `json:"id"`
	Item   string `json:"item"`
	Status string `json:"status"`
}

var (
	orders  = map[string]Order{}
	orderID int64
)

func createOrder(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Item string `json:"item"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}
	if input.Item == "" {
		http.Error(w, `{"error":"item is required"}`, http.StatusUnprocessableEntity)
		return
	}

	id := fmt.Sprintf("ord-%d", atomic.AddInt64(&orderID, 1))
	order := Order{ID: id, Item: input.Item, Status: "created"}
	orders[id] = order

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func getOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	order, ok := orders[id]
	if !ok {
		http.Error(w, `{"error":"order not found"}`, http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func cancelOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	order, ok := orders[id]
	if !ok {
		http.Error(w, `{"error":"order not found"}`, http.StatusNotFound)
		return
	}
	if order.Status == "cancelled" {
		http.Error(w, `{"error":"order already cancelled"}`, http.StatusConflict)
		return
	}

	order.Status = "cancelled"
	orders[id] = order

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /orders", createOrder)
	mux.HandleFunc("GET /orders/{id}", getOrder)
	mux.HandleFunc("PUT /orders/{id}/cancel", cancelOrder)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
