# Section 9: Concurrency

## Beginner → Expert Mapping

| Topic | Level | Importance | Engineering Concept |
|-------|-------|------------|---------------------|
| Goroutines | Beginner | High | Lightweight OS threads, M:N scheduler |
| WaitGroups | Intermediate | High | Barrier synchronization |
| Channels | Advanced | **Critical** | Communicating Sequential Processes (CSP) |
| Select & Context | Advanced | **Critical** | Multiplexing, Cancellation, Timeouts |
| Sync Primitives | Expert | Medium | Mutexes, concurrent map scaling |

## Engineering Depth
Concurrency is Go's flagship feature. Go uses an M:N scheduler, multiplexing thousands of lightweight goroutines (starting at ~2KB memory) onto a handful of OS threads. 
- **Channels:** Unbuffered channels are $O(1)$ synchronization points that *block* until both sender and receiver are ready. Buffered channels behave like blocking queues.
- **Data Races:** Always compile with `-race` (`go test -race ./...`). If two goroutines access the same memory concurrently and at least one is a write, the program is critically compromised.

## References
1. **[Go Blog]** [Share Memory By Communicating](https://go.dev/blog/codelab-share)
2. **[Go Blog]** [Go Concurrency Patterns: Context](https://go.dev/blog/context)
3. **[Effective Go]** [Concurrency](https://go.dev/doc/effective_go#concurrency)

---

## 🏗 Exercise: Concurrent File Downloader (`7-downloader`)

This capstone teaches you the limits of goroutines when dealing with I/O and how to synchronize results.

### Step-by-Step Instructions & Hints
1. **Define URL Target:** Create a slice of dummy API URLs.
2. **Use an Unbuffered Channel:** Create a `results := make(chan string)` to capture download status.
3. **Launch Goroutines:** Loop through the slice and spawn a generic func: `go download(url, results)`.
4. **Wait for Results:** In your `main()` thread, loop the exact number of URLs to pull from the channel: `<-results`.
   - *Hint:* If you loop infinitely, you will hit a `fatal error: all goroutines are asleep - deadlock!`. Channel receivers block until data exists!


## Learning Path

| ID | Lesson | Concept | Requires |
| --- | --- | --- | --- |
| GC.1 | [goroutines](./1-goroutine) | go keyword · M:N scheduler · 2KB stack · closure-capture bug | 🟢 entry |
| GC.2 | [WaitGroups](./2-wait-group) | Add · defer Done · Wait · pass by pointer rule | GC.1 |
| GC.3 | [channels (unbuffered)](./3-channels) | make(chan T) · send · receive · block-until-ready · chan&lt;- / &lt;-chan | GC.1, GC.2 |
| GC.4 | [buffered channels](./4-channels-buffered) | make(chan T, N) · len vs cap · producer-consumer decoupling | GC.3 |
| GC.5 | [closing channels](./5-channels-closing) | close() · range over channel · comma-ok · broadcast signal | GC.3, GC.4 |
| GC.6 | [pipeline project](./6-project-1) | ping/pong actors · select on done channel · context cancel | GC.3, GC.4, GC.5 |
| **GC.8** ⭐ | [race conditions](./8-race) | sync.Mutex · sync.RWMutex · sync/atomic · -race flag | GC.1, GC.3 |
| GC.9 | [select deep dive](./9-select-deep-dive) | Multiplexing · timeout with time.After · non-blocking default · fan-in | GC.3, GC.4, GC.5 |
| GC.10 | [sync primitives](./10-sync-primitives) | sync.Once singleton · sync.Map · when to use each | GC.8 |
