package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type LayerInfo struct {
	Command   string `json:"command"`
	SizeBytes int64  `json:"size_bytes"`
	Cacheable bool   `json:"cacheable"`
}

var recommendedLayers = []LayerInfo{
	{Command: "FROM golang:1.25-alpine", SizeBytes: 150_000_000, Cacheable: true},
	{Command: "WORKDIR /app", SizeBytes: 0, Cacheable: true},
	{Command: "COPY go.mod go.sum ./", SizeBytes: 50_000, Cacheable: true},
	{Command: "RUN go mod download", SizeBytes: 200_000_000, Cacheable: true},
	{Command: "COPY . .", SizeBytes: 500_000, Cacheable: false},
	{Command: "RUN CGO_ENABLED=0 go build -o server .", SizeBytes: 15_000_000, Cacheable: false},
}

type BuildSuggestion struct {
	Layers        []LayerInfo `json:"layers"`
	EstimatedSize int64       `json:"estimated_size_bytes"`
}

func layersHandler(w http.ResponseWriter, r *http.Request) {
	var total int64
	for _, l := range recommendedLayers {
		total += l.SizeBytes
	}
	suggestion := BuildSuggestion{
		Layers:        recommendedLayers,
		EstimatedSize: total,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(suggestion)
}

func main() {
	http.HandleFunc("/layers", layersHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Listening on :%s\n", port)
	http.ListenAndServe(":"+port, nil)
}
