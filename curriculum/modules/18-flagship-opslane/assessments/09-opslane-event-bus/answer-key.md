# Opslane Milestone 09 Assessment — Answer Key

## Multiple Choice (1 point each)

**1. a) To prevent deadlocks and lock contention: if a handler executes slowly or attempts to publish another event, holding the lock would freeze the entire event bus.**  
*Explanation:* Invoking user code while holding a mutex lock is an anti-pattern. If a listener takes 500ms to process an event or invokes `bus.Publish()` internally, holding the lock causes all other publishers to block indefinitely and deadlocks re-entrant calls. Releasing the lock before invocation eliminates this danger.

**2. b) To isolate panics: an unhandled panic in one listener callback will be intercepted and logged, preventing the entire Go runtime process from crashing.**  
*Explanation:* In Go, an unrecovered panic in *any* goroutine crashes the entire process. Wrapping each handler in `defer recover()` confines errors to that specific listener, allowing other listeners and the parent server to keep running.

**3. b) `"*"`.**  
*Explanation:* The bus reserves `"*"` as a wildcard subscription topic. When any event is published, the bus queries both `b.subscribers[e.Type()]` and `b.subscribers["*"]`, dispatching to both.

**4. b) It marks the bus as closed to reject new events, waits for all active handler goroutines to finish via `b.wg.Wait()`, and returns `ctx.Err()` if the context timeout expires first.**  
*Explanation:* Proper graceful shutdown involves two phases: stop accepting new work, and give active work a bounded grace period to complete. Running `b.wg.Wait()` inside a select with `<-ctx.Done()` enforces this pattern cleanly.

**5. b) To allow callers to cleanly `Unsubscribe(topic, subID)` later, avoiding memory leaks when components are destroyed.**  
*Explanation:* Returning a subscription ID gives callers an unforgeable handle to unsubscribe later, preventing memory leaks caused by lingering function closures in long-running processes.

**6. b) `Event` (`Type() string`, `TenantID() string`, `Payload() any`, `OccurredAt() time.Time`).**  
*Explanation:* The `Event` interface establishes the common contract required for topic routing, multi-tenant auditing, and event deserialization.

---

## Short Answer (2 points each)

**7. Synchronous vs. Asynchronous event buses:**
- **Synchronous Bus:** Dispatches handlers sequentially on the publisher's goroutine.  
  *Pros:* Handlers participate in the publisher's transaction/error flow; simple execution order.  
  *Cons:* Publisher is blocked until all handlers finish; slow listeners degrade user-facing API latency; failure in one handler can abort the entire request.
- **Asynchronous Bus:** Dispatches handlers concurrently in background goroutines.  
  *Pros:* Non-blocking for the publisher; slow operations (email, analytics) don't degrade API latency; failures are isolated.  
  *Cons:* Eventual consistency; cannot roll back database transaction if an async handler fails.
*Grading: 1 pt for latency/throughput differences; 1 pt for transactional/consistency trade-offs.*

**8. Goroutine leaks and context cancellation:**
- **What is a Leak:** A goroutine spawned without a guaranteed termination condition (e.g. waiting indefinitely on a channel receive or un-timeouted HTTP call). Leaked goroutines remain in memory forever, consuming stack space and eventually causing OOM crashes.
- **Mitigation:** Passing `ctx context.Context` into event handlers ensures that if an operation takes too long or the server shuts down, the context's cancellation channel (`<-ctx.Done()`) fires. The listener detects this signal, aborts early, and terminates the goroutine.
*Grading: 1 pt for defining goroutine leaks; 1 pt for explaining termination via `ctx.Done()`.*

**9. In-memory event loss and Transactional Outbox Pattern:**
- **In-Memory Vulnerability:** If a service updates the database and publishes an event to an in-memory bus, a crash or power outage immediately after the database commit will wipe all queued in-memory events before they are processed.
- **Transactional Outbox:** The application writes the event to an `outbox` table in the database *within the same transaction* as the order mutation. A separate background worker reads from the outbox table and publishes events to the bus or message broker. If the server crashes, the event remains safely persisted in the database and is processed upon restart.
*Grading: 1 pt for identifying volatile RAM loss on crash; 1 pt for explaining atomic outbox table insertion.*

**10. Coordinated goroutine termination via `sync.WaitGroup`:**
1. **`wg.Add(1)`:** Called on the coordinating goroutine *before* launching the worker goroutine to register pending work.
2. **`defer wg.Done()`:** Executed inside the spawned goroutine upon completion, decrementing the internal counter.
3. **`wg.Wait()`:** Blocks execution until the internal counter reaches 0 (all spawned workers have called `Done()`), ensuring that `Close()` does not exit while background work is active.
*Grading: 1 pt for explaining Add before spawn and Done on exit; 1 pt for explaining Wait blocking until counter reaches zero.*
