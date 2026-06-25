package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
)

type LeakHolder struct {
	data [][]byte
}

func (lh *LeakHolder) leak(size int) {
	buf := make([]byte, size)
	lh.data = append(lh.data, buf)
}

func (lh *LeakHolder) noLeak(size int) {
	_ = make([]byte, size)
}

func main() {
	var holder LeakHolder

	for i := 0; i < 100; i++ {
		holder.leak(1 << 20)
	}

	runtime.GC()

	f, err := os.Create("heap_profile.pprof")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := pprof.WriteHeapProfile(f); err != nil {
		panic(err)
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("HeapAlloc: %d MB\n", m.HeapAlloc/1024/1024)
	fmt.Printf("Heap profile written to heap_profile.pprof\n")

	_ = holder
}
