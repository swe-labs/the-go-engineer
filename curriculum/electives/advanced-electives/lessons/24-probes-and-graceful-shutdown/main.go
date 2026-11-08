package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

const (
	defaultAddr        = ":8080"
	defaultDrainWindow = 10 * time.Second
)

// Server owns the readiness state and the HTTP listener.
type Server struct {
	ready    atomic.Bool
	inflight sync.WaitGroup
	work     http.HandlerFunc
	srv      *http.Server
	ln       net.Listener
}

// NewServer returns a server that is ready. If work is nil, /work sleeps
// briefly and returns 200.
func NewServer(work http.HandlerFunc) *Server {
	s := &Server{work: work}
	if s.work == nil {
		s.work = func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(10 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
		}
	}
	s.ready.Store(true)
	return s
}

// Ready reports whether the server should receive traffic.
func (s *Server) Ready() bool { return s.ready.Load() }

// NotReady removes the server from rotation without stopping it.
func (s *Server) NotReady() { s.ready.Store(false) }

// Start listens on addr and serves in the background.
func (s *Server) Start(addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", addr, err)
	}
	s.ln = ln
	s.srv = &http.Server{Handler: s.routes()}
	go func() {
		if err := s.srv.Serve(ln); err != nil && err != http.ErrServerClosed {
			log.Printf("serve: %v", err)
		}
	}()
	return nil
}

// Addr returns the bound address. It is only valid after Start.
func (s *Server) Addr() string {
	if s.ln == nil {
		return ""
	}
	return s.ln.Addr().String()
}

// Shutdown marks the server not ready, stops accepting new connections, and
// waits for in-flight requests to finish or ctx to expire.
func (s *Server) Shutdown(ctx context.Context) error {
	s.ready.Store(false)
	if s.srv == nil {
		return nil
	}
	err := s.srv.Shutdown(ctx)
	s.inflight.Wait()
	return err
}

func (s *Server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok\n"))
	})
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		if !s.ready.Load() {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte("draining\n"))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ready\n"))
	})
	mux.HandleFunc("/work", s.work)
	return s.track(mux)
}

// track counts in-flight requests so the shutdown path can assert that all
// work finished.
func (s *Server) track(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.inflight.Add(1)
		defer s.inflight.Done()
		next.ServeHTTP(w, r)
	})
}

func main() {
	addr := defaultAddr
	if v := os.Getenv("ADDR"); v != "" {
		addr = v
	}
	s := NewServer(nil)
	if err := s.Start(addr); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("listening on %s; ready=%t\n", s.Addr(), s.Ready())

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sig := <-stop
	fmt.Printf("received %s; draining for up to %s\n", sig, defaultDrainWindow)

	ctx, cancel := context.WithTimeout(context.Background(), defaultDrainWindow)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown: %v", err)
	}
	fmt.Println("shutdown complete")
}
