package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func startServer(addr string) (*http.Server, net.Addr, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, nil, err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	srv := &http.Server{Handler: mux}
	go func() {
		if err := srv.Serve(listener); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
		}
	}()
	return srv, listener.Addr(), nil
}

func waitForSignal(signals ...os.Signal) os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, signals...)
	return <-ch
}

func main() {
	srv, addr, err := startServer(":8080")
	if err != nil {
		panic(err)
	}
	fmt.Printf("server started on %s, waiting for signal\n", addr)

	sig := waitForSignal(syscall.SIGINT, syscall.SIGTERM)
	slog.Info("shutting down", "signal", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("shutdown error", "error", err)
	}
	fmt.Println("server exited gracefully")
}
