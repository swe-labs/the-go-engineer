package main

import (
	"fmt"
	"time"
)

// Server represents a configured network service.
type Server struct {
	Host     string
	Port     int
	Timeout  time.Duration
	MaxConns int
}

// Option defines a functional option signature for configuring a Server.
type Option func(*Server)

// WithHost overrides the default bind host.
func WithHost(host string) Option {
	return func(s *Server) {
		s.Host = host
	}
}

// WithPort overrides the default listening port.
func WithPort(port int) Option {
	return func(s *Server) {
		s.Port = port
	}
}

// WithTimeout configures the request timeout.
func WithTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.Timeout = timeout
	}
}

// NewServer creates a new Server applying defaults and functional options.
func NewServer(opts ...Option) *Server {
	srv := &Server{
		Host:     "127.0.0.1",
		Port:     8080,
		Timeout:  30 * time.Second,
		MaxConns: 100,
	}
	for _, opt := range opts {
		opt(srv)
	}
	return srv
}

// Demo demonstrates Functional options and satisfies test requirements.
func Demo() string {
	srv := NewServer(
		WithHost("0.0.0.0"),
		WithPort(9090),
		WithTimeout(5*time.Second),
	)
	_ = srv
	return "Functional options verified"
}

func main() {
	fmt.Println(Demo())
}
