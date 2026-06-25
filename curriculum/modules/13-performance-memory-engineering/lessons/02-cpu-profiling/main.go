package main

import (
	"fmt"
	"math"
	"os"
	"runtime/pprof"
)

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	for i := 2; i <= int(math.Sqrt(float64(n))); i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func countPrimes(limit int) int {
	count := 0
	for i := 0; i < limit; i++ {
		if isPrime(i) {
			count++
		}
	}
	return count
}

func main() {
	f, err := os.Create("cpu_profile.pprof")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := pprof.StartCPUProfile(f); err != nil {
		panic(err)
	}
	defer pprof.StopCPUProfile()

	result := countPrimes(100000)
	fmt.Printf("Primes up to 100000: %d\n", result)
}
