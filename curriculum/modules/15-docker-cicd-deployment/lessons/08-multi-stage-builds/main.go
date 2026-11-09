package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
)

type BuildInfo struct {
	Version     string `json:"version"`
	GoVersion   string `json:"go_version"`
	GoArch      string `json:"go_arch"`
	GoOS        string `json:"go_os"`
	CgoEnabled  bool   `json:"cgo_enabled"`
	StaticBuild bool   `json:"static_build"`
}

func buildInfoHandler(w http.ResponseWriter, r *http.Request) {
	info := BuildInfo{
		Version:     "1.0.0",
		GoVersion:   runtime.Version(),
		GoArch:      runtime.GOARCH,
		GoOS:        runtime.GOOS,
		CgoEnabled:  false,
		StaticBuild: true,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

func main() {
	http.HandleFunc("/build", buildInfoHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server starting on :%s\n", port)
	http.ListenAndServe(":"+port, nil)
}
