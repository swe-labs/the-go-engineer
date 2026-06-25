package main

import (
	"fmt"
)

type Config struct {
	Port    int
	Timeout int
	Verbose bool
}

func MagicServer() {
	cfg := Config{}
	fmt.Println("MagicServer started on port", cfg.Port)
	fmt.Println("Timeout:", cfg.Timeout, "seconds")
	fmt.Println("Verbose:", cfg.Verbose)
}

func ZeroMagicServer(cfg Config) {
	if cfg.Port == 0 {
		cfg.Port = 8080
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 30
	}
	fmt.Println("ZeroMagicServer started on port", cfg.Port)
	fmt.Println("Timeout:", cfg.Timeout, "seconds")
	fmt.Println("Verbose:", cfg.Verbose)
}

func main() {
	fmt.Println("=== Magic approach (hidden defaults) ===")
	MagicServer()

	fmt.Println()
	fmt.Println("=== Zero-magic approach (explicit defaults) ===")
	ZeroMagicServer(Config{Port: 9090, Timeout: 60, Verbose: true})
	ZeroMagicServer(Config{})
}
