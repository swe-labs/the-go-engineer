package main

import (
	"fmt"
	"net/http"
	"os"
)

// TODO: Replace this scaffold with your Module 15 service wired for Kubernetes.
//
// Required endpoints:
//   - GET /healthz : liveness. Return 200 while the process is alive.
//   - GET /readyz  : readiness. Return 200 only when the process can serve
//     traffic, and 503 during startup and shutdown.
//   - GET /version : return the build revision and environment as JSON so a
//     rollout can prove which revision answered.
//
// Required behavior:
//   - Read configuration from the environment and from a mounted Secret file.
//   - Validate configuration at startup and exit non-zero on a missing value.
//   - On SIGTERM, flip readiness to false, stop accepting connections, drain
//     in-flight requests within the termination grace period, then exit 0.
//   - Set GOMEMLIMIT from the container memory limit and report GOMAXPROCS.

// TODO: implement the readiness state so /readyz reflects it.
type Server struct{}

// TODO: implement the routes. Register /healthz, /readyz, and /version.
func (s *Server) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	return mux
}

func main() {
	// TODO: load and validate configuration before serving.
	// TODO: create the listener, handle SIGTERM, and drain gracefully.
	addr := ":8080"
	if v := os.Getenv("APP_ADDR"); v != "" {
		addr = v
	}
	fmt.Printf("listening on %s\n", addr)

	s := &Server{}
	if err := http.ListenAndServe(addr, s.routes()); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}
