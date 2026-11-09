package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type ReportHeader struct {
	Type      string `json:"type"`
	Timestamp string `json:"timestamp"`
	Total     int    `json:"total"`
}

type ReportItem struct {
	ID    int    `json:"id"`
	Value string `json:"value"`
}

type ReportFooter struct {
	Status string `json:"status"`
}

func reportHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Stream the JSON structure: header, array of items, footer.
	encoder := json.NewEncoder(w)
	encoder.Encode(map[string]string{"header": "report"})
	encoder.Encode(ReportHeader{
		Type:      "summary",
		Timestamp: time.Now().Format(time.RFC3339),
		Total:     100,
	})
	encoder.Encode(map[string]string{"items": "start"})
	for i := 0; i < 100; i++ {
		encoder.Encode(ReportItem{
			ID:    i + 1,
			Value: fmt.Sprintf("item-%d", i+1),
		})
	}
	encoder.Encode(map[string]string{"items": "end"})
	encoder.Encode(ReportFooter{Status: "complete"})
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"message": "pong"})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /ping", pingHandler)
	mux.HandleFunc("GET /report", reportHandler)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
