# HS.7 Server Timeouts

## Mission

Protect your server from resource exhaustion by implementing production-grade timeouts for network I/O and request processing.

## Prerequisites

- `HS.6` error-handling-middleware

## Mental Model

Think of timeouts as **Table Turn Limits at a Busy Restaurant**.

1. **Arrival Timeout (`ReadTimeout`)**: If a guest sits down but takes 20 minutes to even look at the menu (sending the request headers), the waiter asks them to leave to free up the table.
2. **Dining Timeout (`TimeoutHandler`)**: The chef gives each order a 10-minute limit. If the food isn't ready by then, the kitchen sends an apology (503 Service Unavailable) instead of making the guest wait forever.
3. **Departure Timeout (`WriteTimeout`)**: Once the food is served, if the guest takes an hour to eat a single fry (slow network read), the table is cleared.
4. **The Idle Table (`IdleTimeout`)**: How long can a guest stay at a table after they've finished eating before we ask them to leave or order something else?

## Visual Model

```mermaid
graph LR
    A["Request Start"] -- "ReadTimeout" --> B["Headers/Body Read"]
    B -- "TimeoutHandler" --> C["Business Logic"]
    C -- "WriteTimeout" --> D["Response Sent"]
    D -- "IdleTimeout" --> E["Connection Re-use"]
```

## Machine View

Every open connection on your server consumes a file descriptor and memory (for the goroutine stack and buffers). Without timeouts, a malicious or broken client can open thousands of connections and simply "stay quiet." This is the basis of a **Slowloris DoS attack**. By setting `ReadTimeout` and `WriteTimeout` on the `http.Server` struct, you instruct the Go runtime to automatically close any TCP connection that stays open longer than allowed without making progress. `http.TimeoutHandler` is different: it works at the logic level, returning a 503 error if your code doesn't finish within a certain time, preventing a hanging database call from tying up a worker goroutine forever.

## Run Instructions

```bash
go run ./06-backend-db/01-web-and-database/http-servers/7-server-timeouts
```

Observe the timeout in action:
```bash
# This will succeed instantly
curl -i http://localhost:8086/fast

# This will wait 2 seconds and then return a 503 Error
curl -i http://localhost:8086/slow
```

## Code Walkthrough

### `http.Server` Configuration
Instead of using `http.ListenAndServe`, we create an instance of `http.Server`. This is where we define the global rules for our network listener.

### `ReadTimeout` vs `WriteTimeout`
- **Read**: Starts from when the connection is accepted until the request body is fully read.
- **Write**: Starts from the end of the request header read until the end of the response write.

### `http.TimeoutHandler`
This middleware allows you to set a per-route or global deadline for your logic. It is safer than just using a context timeout because it also handles the HTTP response part (sending a 503) for you.

### Why default to 503?
When a timeout occurs, the server cannot fulfill the request. `503 Service Unavailable` is the standard way to tell the client: "I'm too busy or taking too long, please try again later."

## Try It

1. Change the `TimeoutHandler` limit to be longer than the `time.Sleep` in the handler and verify the request succeeds.
2. Set a very small `ReadHeaderTimeout` (e.g., `100ms`) and try to connect using a slow tool or manual telnet.
3. Observe what happens to the console logs when a `TimeoutHandler` triggers-does the slow logic stop immediately? (Hint: See the Context lesson in Section 07).

## In Production
Timeouts are not "One Size Fits All". A file upload service might need a 10-minute `ReadTimeout`, while a high-frequency trading API might need a 10ms timeout. Always monitor your "p99" response times and set your timeouts slightly above your expected maximum latency to provide a safety margin without cutting off legitimate users.

## Thinking Questions
1. What is the difference between a network timeout and a logic timeout?
2. How can a single hanging request cause a "Cascading Failure" in a microservice architecture?
3. Why does Go's `http.ListenAndServe` not include default timeouts?

> **Forward Reference:** You've secured your server's resources. Now, how do you shut it down without cutting off active users? In [Lesson 8: Graceful HTTP Shutdown](../8-graceful-http-shutdown/README.md), you will learn how to exit cleanly and finish pending work.

## Next Step

Next: `HS.8` -> `06-backend-db/01-web-and-database/http-servers/8-graceful-http-shutdown`

Open `06-backend-db/01-web-and-database/http-servers/8-graceful-http-shutdown/README.md` to continue.
