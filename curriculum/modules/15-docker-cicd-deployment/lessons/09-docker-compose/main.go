package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type AppConfig struct {
	ServiceName string `json:"service_name"`
	Version     string `json:"version"`
	DatabaseURL string `json:"database_url"`
	RedisURL    string `json:"redis_url"`
	ListenAddr  string `json:"listen_addr"`
}

func configHandler(w http.ResponseWriter, r *http.Request) {
	cfg := AppConfig{
		ServiceName: getEnv("SERVICE_NAME", "api"),
		Version:     getEnv("VERSION", "1.0.0"),
		DatabaseURL: getEnv("DATABASE_URL", "not-configured"),
		RedisURL:    getEnv("REDIS_URL", "not-configured"),
		ListenAddr:  getEnv("LISTEN_ADDR", ":8080"),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cfg)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	http.HandleFunc("/config", configHandler)
	http.HandleFunc("/health", healthHandler)

	addr := getEnv("LISTEN_ADDR", ":8080")
	fmt.Printf("Server starting on %s\n", addr)
	http.ListenAndServe(addr, nil)
}
