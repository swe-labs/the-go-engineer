package main

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"
)

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// recoveryWriter wraps ResponseWriter to track whether headers were written.
type recoveryWriter struct {
	http.ResponseWriter
	wroteHeader bool
	status      int
}

func (w *recoveryWriter) WriteHeader(code int) {
	if !w.wroteHeader {
		w.status = code
		w.wroteHeader = true
		w.ResponseWriter.WriteHeader(code)
	}
}

func (w *recoveryWriter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(b)
}

func withRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &recoveryWriter{ResponseWriter: w, status: http.StatusOK}
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("PANIC: %v\n%s", rec, debug.Stack())
				if !rw.wroteHeader {
					writeJSON(rw, http.StatusInternalServerError, APIError{
						Code:    "PANIC",
						Message: "internal error",
					})
				}
			}
		}()
		next.ServeHTTP(rw, r)

		// If the handler wrote a 5xx without a proper JSON body, wrap it.
		if rw.status >= 500 && !rw.wroteHeader {
			// Headers already sent; we can't replace the body.
			// In a real app, the handler should use writeJSON.
			log.Printf("WARNING: handler returned %d without JSON body", rw.status)
		}
	})
}

func okHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func panicHandler(w http.ResponseWriter, r *http.Request) {
	panic("simulated database failure")
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusNotFound, APIError{
		Code:    "NOT_FOUND",
		Message: "resource not found",
	})
}

func plainErrorHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "something went wrong", http.StatusInternalServerError)
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/ok", withRecovery(http.HandlerFunc(okHandler)))
	mux.Handle("/panic", withRecovery(http.HandlerFunc(panicHandler)))
	mux.Handle("/notfound", withRecovery(http.HandlerFunc(notFoundHandler)))
	mux.Handle("/plain-error", withRecovery(http.HandlerFunc(plainErrorHandler)))

	log.Fatal(http.ListenAndServe(":8080", mux))
}
