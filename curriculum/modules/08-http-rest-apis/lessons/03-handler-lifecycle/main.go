package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type Logger struct {
	handler http.Handler
}

func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	l.handler.ServeHTTP(w, r)
	log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
}

type AdminCheck struct {
	handler http.Handler
}

func (a *AdminCheck) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Admin") != "true" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	a.handler.ServeHTTP(w, r)
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Welcome, admin!")
}

func main() {
	var h http.Handler = http.HandlerFunc(adminHandler)
	h = &AdminCheck{handler: h}
	h = &Logger{handler: h}

	http.Handle("/admin", h)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
