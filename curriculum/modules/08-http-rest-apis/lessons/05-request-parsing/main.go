package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type InputData struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type ErrorResponse struct {
	Errors map[string]string `json:"errors"`
}

func processHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeErrors(w, http.StatusBadRequest, map[string]string{"id": "path parameter {id} is required"})
		return
	}

	format := r.URL.Query().Get("format")
	if format == "" {
		format = "json"
	}
	if format != "json" {
		writeErrors(w, http.StatusBadRequest, map[string]string{"format": "unsupported format, use 'json'"})
		return
	}

	traceID := r.Header.Get("X-Trace-ID")

	var input InputData
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeErrors(w, http.StatusBadRequest, map[string]string{"body": "invalid JSON: " + err.Error()})
		return
	}

	errs := make(map[string]string)
	if input.Name == "" {
		errs["name"] = "field 'name' is required"
	}
	if input.Count == 0 {
		errs["count"] = "field 'count' is required and must be non-zero"
	}
	if len(errs) > 0 {
		writeErrors(w, http.StatusBadRequest, errs)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":       id,
		"name":     input.Name,
		"count":    input.Count,
		"format":   format,
		"trace_id": traceID,
	})
}

func writeErrors(w http.ResponseWriter, status int, errs map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Errors: errs})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /process/{id}", processHandler)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
