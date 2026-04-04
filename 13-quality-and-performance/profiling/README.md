# Section 25: Profiling with pprof

## Beginner → Expert Mapping

| Topic | Level | Importance | Engineering Concept |
|-------|-------|------------|---------------------|
| CPU profiling | Intermediate | **Critical** | Flame graphs, hot path identification |
| Memory profiling | Intermediate | **Critical** | Heap allocation, inuse vs alloc |
| `net/http/pprof` | Advanced | High | Live production profiling endpoint |
| `runtime/trace` | Expert | High | Goroutine scheduling, GC events |

## Engineering Depth

`pprof` answers questions no benchmark can: "Where does my service spend its time when serving real production traffic?" A benchmark tells you one function is 200 ns/op. pprof tells you that function runs 200,000 times per second and accounts for 40% of total CPU.

**The profiling workflow:**
1. Record a profile (`-cpuprofile` flag or `pprof.StartCPUProfile`)
2. Visualise with `go tool pprof` — text, web, or flame graph
3. Fix the bottleneck (usually one of: reflection, string building, excessive allocation, regex in a hot loop)
4. Benchmark before and after to confirm the improvement
5. Commit the benchmark so regressions surface in CI

**Common hotspots found via pprof:**
- `runtime.mallocgc` — too many allocations (use sync.Pool or pre-allocate)
- `runtime.gcBgMarkWorker` — GC running too frequently (same root cause)
- `regexp.(*Regexp).FindAllString` — compiling regex inside a loop (pre-compile at package level)
- `strings.Builder.copyCheck` — using += in a loop (use strings.Builder)
- `reflect.Value.Field` — encoding/json reflection on untagged large structs

## How to Run

```bash
# CPU profile
go run ./25-profiling/1-cpu-profile
go tool pprof -http=:8090 cpu.prof

# Memory profile
go run ./25-profiling/2-memory-profile
go tool pprof -http=:8090 mem.prof

# Live pprof endpoint
go run ./25-profiling/3-http-pprof
# Then: go tool pprof http://localhost:8080/debug/pprof/profile?seconds=5
```

## References

- [Go Blog: Profiling Go Programs](https://go.dev/blog/pprof)
- [Package pprof](https://pkg.go.dev/net/http/pprof)
- [Package runtime/pprof](https://pkg.go.dev/runtime/pprof)
- [Package runtime/trace](https://pkg.go.dev/runtime/trace)
