# Incident debugging

## Learning objective

Apply hypothesis-driven debugging methodology to production Go incidents, using structured observability data, pprof profiling, and systematic root-cause analysis.

## Why this matters

Production incidents are inevitable. The difference between a 10-minute outage and a 2-hour outage is almost never technical skill — it is methodology. Engineers who jump to conclusions (restart the pod, scale up, change config) without gathering data often make things worse. Engineers who follow a structured process — observe, correlate, isolate, diagnose, remediate — resolve incidents faster and prevent recurrence. This lesson teaches the process that SREs at Google, Stripe, and GitHub use every day.

## Mental model

Incident debugging is a scientific investigation, not a guessing game. The process mirrors the scientific method:

1. **Observe**: What does the data say? Check dashboards, alerts, logs. Define the blast radius.
2. **Hypothesize**: Form a falsifiable hypothesis. "The database connection pool is exhausted because the query cache was invalidated."
3. **Predict**: If my hypothesis is correct, I should see X. "I should see 100+ connections in `pg_stat_activity` from this service."
4. **Test**: Collect the predicted evidence. Run `SELECT * FROM pg_stat_activity WHERE application_name = 'orders'`.
5. **Conclude**: Hypothesis supported or refuted. If refuted, return to step 2.
6. **Remediate**: Fix root cause, not symptom. Then verify the fix.

The analogy breaks because production incidents have time pressure — you cannot spend hours testing every hypothesis. The skill is prioritizing the highest-probability hypothesis based on available data.

## Core idea

Structured incident response follows a defined severity-based process:

| Phase | Action | Tools |
|---|---|---|
| Detection | Alert fires or user reports issue | PagerDuty, Prometheus, Grafana |
| Triage | Determine severity and blast radius | Dashboard, logs, traces |
| Diagnosis | Find root cause using hypotheses | pprof, tracing, metrics, logging |
| Mitigation | Stop the bleeding | Rollback, scale out, feature flag |
| Resolution | Apply permanent fix | Code change, config change |
| Review | Post-incident review | Blameless postmortem |

Key debugging tools in Go:

- **pprof**: CPU, heap, goroutine, mutex, and block profiles via `net/http/pprof`.
- **GODEBUG** environment variable: `gctrace=1` for GC log, `schedtrace=100` for scheduler trace.
- **tracing**: `go tool trace` for execution trace (goroutine scheduling, GC, network events).
- **strace/dlv**: System call tracing and Go debugger for local reproduction.

## Under the hood

The Go `runtime/pprof` package samples the call stack at a configurable interval. For CPU profiling, the runtime's SIGPROF handler interrupts execution every 10ms (default) and records the program counter. After collection, `pprof` uses the compiled binary's symbol table to map PCs to function names and line numbers. The profile is a set of stack traces with counts, represented as a weighted directed graph.

`net/http/pprof` registers handlers on the default mux:
- `/debug/pprof/goroutine`: stack traces of all goroutines.
- `/debug/pprof/heap`: heap memory allocations.
- `/debug/pprof/profile`: CPU profile (30s default).
- `/debug/pprof/mutex`: contended mutexes.
- `/debug/pprof/block`: blocking operations.

The heap profile uses _sampling_: it records every allocation over 512KB (default `MemProfileRate`). It is not exact but statistically represents allocation patterns. The goroutine profile captures the current state of every goroutine, including the function that started it and its current blocking point.

## How Go uses it

Production Go services always enable the pprof HTTP endpoints (but restricted to internal networks). When an incident occurs:

1. Fetch goroutine profile: `curl http://service:8080/debug/pprof/goroutine?debug=2`
2. Look for goroutines stuck in `Gwaiting` (waiting on channel, mutex, network).
3. Fetch heap profile: `curl http://service:8080/debug/pprof/heap?debug=1`
4. Compare with baseline to find memory leaks.
5. Fetch CPU profile: `go tool pprof http://service:8080/debug/pprof/profile?seconds=30`
6. Analyze with `top`, `web`, `peek` commands.

## Go example

```go
package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"sync"
	"time"
)

var (
	leakedGoroutines sync.WaitGroup
	mu               sync.Mutex
	contendedCounter int
)

func main() {
	http.HandleFunc("/api/leak", leakHandler)
	http.HandleFunc("/api/work", workHandler)
	http.HandleFunc("/api/contend", contendHandler)

	go func() {
		log.Println("pprof available at /debug/pprof/")
		log.Println("Simulate: open http://localhost:8080/debug/pprof/")
		fmt.Println("\nReady. Trigger scenarios:")
		fmt.Println("  GET /api/leak    — leaks goroutines")
		fmt.Println("  GET /api/work    — CPU-intensive work")
		fmt.Println("  GET /api/contend — mutex contention")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	select {}
}

func leakHandler(w http.ResponseWriter, r *http.Request) {
	for i := 0; i < 10; i++ {
		leakedGoroutines.Add(1)
		go func(id int) {
			defer leakedGoroutines.Done()
			ch := make(chan struct{})
			<-ch // blocks forever — leak
		}(i)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"leaked": 10}`))
}

func workHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	result := 0
	for i := 0; i < 10_000_000; i++ {
		result += rand.Intn(1000)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"result": %d, "duration": "%v"}`, result, time.Since(start))
}

func contendHandler(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
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
```

## Step-by-step execution

1. Start the server. Open `http://localhost:8080/debug/pprof/` in browser — shows available profiles.
2. Call `GET /api/leak` — spawns 10 goroutines blocked on channel receive (simulated leak).
3. Fetch `http://localhost:8080/debug/pprof/goroutine?debug=2` — see goroutines stuck in `chan receive`.
4. Count goroutines with state `Gwaiting` and stack trace showing `leakHandler`.
5. Call `GET /api/work` — does CPU-bound work.
6. Fetch CPU profile: `go tool pprof http://localhost:8080/debug/pprof/profile?seconds=10`.
7. In pprof, type `top` — see `workHandler` and `rand.Intn` as top consumers.
8. Call `GET /api/contend` — 20 goroutines contending on mutex.
9. Fetch mutex profile: `http://localhost:8080/debug/pprof/mutex?debug=1`.
10. Use collected data to write a structured incident report with timeline, root cause, and action items.

## Common mistakes

- Mistake: Starting an incident investigation by guessing the root cause rather than gathering data.
  - Why it happens: Confirmation bias — the engineer already saw a similar issue and assumes the same cause.
  - Fix: Write down your hypothesis and what data would prove or disprove it before looking at any code.

- Mistake: Jumping straight into code without checking the dashboard first.
  - Why it happens: Engineers are comfortable reading code and default to what they know.
  - Fix: The dashboard tells you the blast radius, affected users, and timeline. Skipping it means investigating blind.

- Mistake: Focusing on a single log line instead of the full trace.
  - Why it happens: A single error log is more salient than the surrounding context.
  - Fix: Always filter logs by correlation ID to see the full sequence of events for the affected request.

- Mistake: Making changes to production (restart, scale) before understanding root cause.
  - Why it happens: Urgency to restore service overrides the need to understand root cause.
  - Fix: Apply a targeted mitigation (rollback bad deploy, scale out) first, then investigate root cause after the service is stable.

## Debugging walkthrough

A Go HTTP service serving 5000 req/s suddenly shows 5-second latency spikes every 30 seconds. CPU is at 30%, memory is stable.

**Symptom**: p99 latency jumps from 100ms to 5s every 30s.

**Investigation**:
1. Check dashboard: latency spikes are periodic, every 30 seconds.
2. Hypothesis: GC pause is causing the latency spike. GC in Go runs every 2 minutes by default, but with GOGC=100 and high allocation rate, it could run every 30s.
3. Predict: The goroutine profile should show STW (stop-the-world) markers, and the heap profile should show high allocation rate.
4. Test: Fetch `/debug/pprof/goroutine?debug=2` during a spike — goroutines show `runtime.gcBgMarkWorker` and `runtime.gcAssistAlloc`.
5. Test: Fetch CPU profile during spike — `mallocgc` and `gcDrain` appear in top.
6. Conclude: Hypothesis supported. GC is running every 30s because of high allocation rate.
7. Remediate: Profile heap to find allocation hot spots, reduce allocations (object pooling, pre-allocation), or adjust GOGC.

**Fix**: The handler was allocating a new `bytes.Buffer` per request. Changed to `sync.Pool` for buffer reuse. GC frequency dropped from every 30s to every 5 minutes.

## Production notes

- Always enable pprof in production, but restrict to internal network or require authentication. Never expose `/debug/pprof/` to the public internet.
- Use `go tool pprof -base` to compare two heap profiles (before and after memory leak) to identify leaked objects.
- Set up continuous profiling: tools like Parca or Pyroscope capture profiles at regular intervals and store them for historical comparison.
- For every incident, write a blameless postmortem within 72 hours. Include timeline, root cause, impact, action items with owners and deadlines.
- Practice incident response with chaos engineering and game days. Run "fire drills" where a team member introduces a failure and another must diagnose it using pprof and metrics.

## Performance implications

- pprof CPU profiling adds ~5% CPU overhead during collection due to the SIGPROF signal handler running every 10ms.
- Heap profiling records every allocation over 512KB — memory overhead is negligible.
- Goroutine profile captures all goroutine stacks — for 100K goroutines, this can take 100ms and allocate 50MB. Use sparingly on high-density services.
- Mutex and block profiles require explicit enabling via `runtime.SetMutexProfileFraction(1)` and `runtime.SetBlockProfileRate(1)`. Profile collection adds ~1% overhead in the instrumented paths.

## Practice task

Write a function `FetchAndAnalyzeProfiles(baseURL string) (*ProfileReport, error)` that:
1. Fetches the goroutine profile from `baseURL/debug/pprof/goroutine?debug=2`.
2. Parses the output to count goroutines in each state (running, waiting, syscall, etc.).
3. Fetches the heap profile from `baseURL/debug/pprof/heap?debug=1`.
4. Parses the output to find the top 3 allocation sites by total bytes.
5. Returns a `ProfileReport` struct with goroutine count, blocked goroutines, and top allocations.

## Tests / verification

```bash
go test ./curriculum/modules/12-observability-diagnostics/lessons/11-incident-debugging -v
go run ./curriculum/modules/12-observability-diagnostics/lessons/11-incident-debugging
```

## Review questions

1. What is the first thing you should do when a production alert fires?
2. Why is it dangerous to restart a pod during incident investigation without capturing a goroutine profile first?
3. What does a heap profile tell you that a goroutine profile does not?
4. How does the scientific method (hypothesize, predict, test) apply to incident debugging?
5. What should a blameless postmortem include and why is "blameless" important?

## NEXT UP

Reliability review — systematic patterns for building resilient Go services.
