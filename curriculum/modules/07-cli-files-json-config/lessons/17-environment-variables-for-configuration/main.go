package main

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port    int
	DBURL   string
	Verbose bool
}

func ConfigFromEnv() Config {
	cfg := Config{
		Port:    8080,
		DBURL:   "postgres://localhost:5432/app?sslmode=disable",
		Verbose: false,
	}
	if v := os.Getenv("PORT"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			cfg.Port = p
		}
	}
	if v := os.Getenv("DATABASE_URL"); v != "" {
		cfg.DBURL = v
	}
	if v, ok := os.LookupEnv("VERBOSE"); ok {
		cfg.Verbose = v == "true"
	}
	return cfg
}

func main() {
	// Before setting env vars
	cfg := ConfigFromEnv()
	fmt.Printf("Default config: Port=%d DBURL=%s Verbose=%v\n", cfg.Port, cfg.DBURL, cfg.Verbose)

	// Set env vars for demonstration
	os.Setenv("PORT", "3000")
	os.Setenv("DATABASE_URL", "postgres://prod:5432/db?sslmode=require")
	os.Setenv("VERBOSE", "true")

	cfg2 := ConfigFromEnv()
	fmt.Printf("After setenv:  Port=%d DBURL=%s Verbose=%v\n", cfg2.Port, cfg2.DBURL, cfg2.Verbose)

	// LookupEnv for missing variable
	val, ok := os.LookupEnv("MISSING_VAR")
	fmt.Printf("LookupEnv(MISSING_VAR): val=%q ok=%v\n", val, ok)
}
