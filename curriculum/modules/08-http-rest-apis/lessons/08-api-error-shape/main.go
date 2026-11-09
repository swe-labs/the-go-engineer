package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type APIError struct {
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Details []FieldError `json:"details,omitempty"`
}

func writeError(w http.ResponseWriter, status int, err APIError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(err)
}

var errNotFound = APIError{
	Code:    "NOT_FOUND",
	Message: "The requested resource was not found",
}

var errInternal = APIError{
	Code:    "INTERNAL_ERROR",
	Message: "An unexpected error occurred",
}

func deleteResourceHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if id == "" {
		writeError(w, http.StatusBadRequest, APIError{
			Code:    "BAD_REQUEST",
			Message: "Missing resource ID",
		})
		return
	}

	if id == "panic" {
		panic("simulated panic for testing")
	}

	if id != "42" {
		writeError(w, http.StatusNotFound, errNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /resources/{id}", deleteResourceHandler)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
