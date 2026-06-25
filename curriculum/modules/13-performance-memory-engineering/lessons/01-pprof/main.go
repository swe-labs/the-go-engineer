package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

func allocateMemory(size int) []byte {
	return make([]byte, size)
}

func main() {
	cpuFile, err := os.Create("cpu.pprof")
	if err != nil {
		panic(err)
	}
	defer cpuFile.Close()

	if err := pprof.StartCPUProfile(cpuFile); err != nil {
		panic(err)
	}
	defer pprof.StopCPUProfile()

	result := fibonacci(40)
	fmt.Printf("fibonacci(40) = %d\n", result)

	_ = allocateMemory(10 << 20)

	heapFile, err := os.Create("heap.pprof")
	if err != nil {
		panic(err)
	}
	defer heapFile.Close()

	runtime.GC()
	if err := pprof.WriteHeapProfile(heapFile); err != nil {
		panic(err)
	}

	fmt.Println("Profiles written to cpu.pprof and heap.pprof")

	time.Sleep(100 * time.Millisecond)
}
