package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func greetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Hello, Gopher!")
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"now": time.Now().Format(time.RFC3339)})
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Go 1.24")
}

func main() {
	http.HandleFunc("/greet", greetHandler)
	http.HandleFunc("/time", timeHandler)
	http.HandleFunc("/version", versionHandler)

	addr := ":8080"
	log.Printf("Server starting on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
