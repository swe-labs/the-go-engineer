package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
)

// Budget describes the CPU and memory budget the process should respect.
type Budget struct {
	GOMAXPROCS  int
	MemoryLimit int64 // bytes used for GOMEMLIMIT
	GCPercent   int
}

// ComputeBudget builds a Budget from the runtime and an explicit container
// memory limit. It reserves headroom for the Go runtime and the page cache,
// then sets GOMEMLIMIT below the container limit so the garbage collector runs
// before the kernel OOM killer does.
func ComputeBudget(containerLimitBytes int64, gcPercent int) (Budget, error) {
	if containerLimitBytes <= 0 {
		return Budget{}, fmt.Errorf("container memory limit %d must be positive", containerLimitBytes)
	}
	if gcPercent < 0 {
		return Budget{}, fmt.Errorf("gc percent %d must not be negative", gcPercent)
	}

	b := Budget{
		GOMAXPROCS:  runtime.GOMAXPROCS(0),
		MemoryLimit: containerLimitBytes * 85 / 100,
		GCPercent:   gcPercent,
	}
	if gcPercent > 0 {
		debug.SetGCPercent(gcPercent)
	}
	debug.SetMemoryLimit(b.MemoryLimit)
	return b, nil
}

// parseGCPercent reads the GOGC environment value. An empty value keeps the Go
// default of 100.
func parseGCPercent(raw string) (int, error) {
	if raw == "" {
		return 100, nil
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < 0 {
		return 0, fmt.Errorf("GOGC %q must be a non-negative integer", raw)
	}
	return n, nil
}

func main() {
	limit := int64(512 << 20)
	if raw := os.Getenv("APP_MEMORY_LIMIT_BYTES"); raw != "" {
		n, err := strconv.ParseInt(raw, 10, 64)
		if err != nil || n <= 0 {
			fmt.Fprintf(os.Stderr, "invalid APP_MEMORY_LIMIT_BYTES %q\n", raw)
			os.Exit(1)
		}
		limit = n
	}
	gc, err := parseGCPercent(os.Getenv("GOGC"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		os.Exit(1)
	}
	b, err := ComputeBudget(limit, gc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		os.Exit(1)
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("gomaxprocs=%d container_limit=%d soft_limit=%d gc_percent=%d heap_alloc=%d\n",
		b.GOMAXPROCS, limit, b.MemoryLimit, b.GCPercent, m.HeapAlloc)
}
