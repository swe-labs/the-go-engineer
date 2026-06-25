package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"runtime"
)

type Info struct {
	Service  string `json:"service"`
	Version  string `json:"version"`
	GoArch   string `json:"go_arch"`
	GoOS     string `json:"go_os"`
	Hostname string `json:"hostname"`
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	info := Info{
		Service:  "docker-basics",
		Version:  "1.0.0",
		GoArch:   runtime.GOARCH,
		GoOS:     runtime.GOOS,
		Hostname: hostname,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

func main() {
	addr := os.Getenv("LISTEN_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	http.HandleFunc("/info", infoHandler)
	slog.Info("server starting", "addr", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
