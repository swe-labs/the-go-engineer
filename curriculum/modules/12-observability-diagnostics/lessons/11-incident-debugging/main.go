package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"sync"
	"time"
)

var (
	mu               sync.Mutex
	contendedCounter int
)

func init() {
	runtime.SetMutexProfileFraction(5)
	runtime.SetBlockProfileRate(1)
}

func main() {
	http.HandleFunc("/api/leak", leakHandler)
	http.HandleFunc("/api/work", workHandler)
	http.HandleFunc("/api/contend", contendHandler)

	go func() {
		log.Println("Server starting on :8080")
		log.Println("Profiles at /debug/pprof/")
		fmt.Println("\nEndpoints:")
		fmt.Println("  GET /api/leak    — simulate goroutine leak")
		fmt.Println("  GET /api/work    — simulate CPU work")
		fmt.Println("  GET /api/contend — simulate mutex contention")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	select {}
}

func leakHandler(w http.ResponseWriter, r *http.Request) {
	count := 10
	for i := 0; i < count; i++ {
		go func(id int) {
			ch := make(chan struct{})
			<-ch
		}(i)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"leaked": %d}`, count)
}

func workHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	result := 0
	for i := 0; i < 5_000_000; i++ {
		result += rand.Intn(1000)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"result": %d, "duration": "%v"}`, result, time.Since(start))
}

func contendHandler(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100_000; j++ {
				mu.Lock()
				contendedCounter++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"counter": %d}`, contendedCounter)
}
