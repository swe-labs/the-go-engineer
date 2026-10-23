# Opslane Milestone 09 Assessment — Event Bus: Questions

Total: 10 questions (6 multiple choice × 1 pt, 4 short answer × 2 pt = 14 points). Passing score: 10/14 (71%).

---

## Multiple Choice (1 point each)

**1.** Why does `EventBus.Publish` copy the subscriber slice and release `b.mu.RUnlock()` *before* invoking handler functions?
- a) To prevent deadlocks and lock contention: if a handler executes slowly or attempts to publish another event, holding the lock would freeze the entire event bus.
- b) Go slices cannot be read by multiple goroutines.
- c) To clear the subscribers map.
- d) Because `RUnlock` is required before calling `go func()`.

**2.** What is the purpose of `defer func() { if r := recover(); r != nil { ... } }()` inside each handler goroutine?
- a) To automatically retry failed database queries.
- b) To isolate panics: an unhandled panic in one listener callback will be intercepted and logged, preventing the entire Go runtime process from crashing.
- c) To restart the operating system.
- d) To send an HTTP 500 response.

**3.** In `EventBus`, what special topic allows a listener to receive all events published across all topics?
- a) `"all"`.
- b) `"*"`.
- c) `""` (empty string).
- d) `"default"`.

**4.** How does `EventBus.Close(ctx)` guarantee graceful draining?
- a) It terminates all goroutines with `os.Exit(0)`.
- b) It marks the bus as closed to reject new events, waits for all active handler goroutines to finish via `b.wg.Wait()`, and returns `ctx.Err()` if the context timeout expires first.
- c) It writes all pending events to a text file.
- d) It reboots the database.

**5.** Why does `Subscribe` return an integer subscription ID?
- a) To bill the subscriber per event received.
- b) To allow callers to cleanly `Unsubscribe(topic, subID)` later, avoiding memory leaks when components are destroyed.
- c) To specify the priority order of the handler.
- d) To limit the number of handlers to 100.

**6.** What interface must all published events satisfy in Opslane's Event Bus?
- a) `io.Reader`.
- b) `Event` (`Type() string`, `TenantID() string`, `Payload() any`, `OccurredAt() time.Time`).
- c) `json.Marshaler`.
- d) `fmt.Stringer`.

---

## Short Answer (2 points each)

**7.** Explain the difference between synchronous in-process event buses and asynchronous event buses. What are the latency and transactional trade-offs?

**8.** What is a goroutine leak? How could an improperly designed event bus listener cause goroutine leaks, and how does `context.Context` cancellation mitigate this?

**9.** Explain the Transactional Outbox Pattern. Why is an in-memory event bus vulnerable to message loss during power failures or process crashes?

**10.** How does `sync.WaitGroup` coordinate asynchronous goroutine termination in `Close()`? Explain the sequence of `wg.Add(1)`, `wg.Done()`, and `wg.Wait()`.
