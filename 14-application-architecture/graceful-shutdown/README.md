# Section 27: Graceful Shutdown

## Beginner → Expert Mapping

| Topic | Level | Importance | Engineering Concept |
|-------|-------|------------|---------------------|
| OS signals | Beginner | **Critical** | SIGTERM, SIGINT, signal.NotifyContext |
| HTTP graceful drain | Intermediate | **Critical** | server.Shutdown, in-flight request handling |
| Complete graceful server | Advanced | **Critical** | Signals + HTTP + DB + background workers |

## Engineering Depth

`log.Fatal(http.ListenAndServe(":8080", mux))` is in every tutorial. It is wrong for production. When Kubernetes sends SIGTERM to terminate a pod, this crashes instantly — killing in-flight requests. Users see 502 errors during every deployment.

**The zero-downtime deployment sequence:**
1. Kubernetes sets pod to "Terminating" and removes it from load balancer (takes ~2-5 seconds to propagate)
2. Kubernetes sends SIGTERM
3. Your server receives SIGTERM and calls `http.Server.Shutdown(ctx)` with a 30-second deadline
4. `Shutdown()` stops accepting NEW connections
5. `Shutdown()` waits for existing connections to drain
6. After drain (or deadline), the process exits cleanly

Without graceful shutdown, step 2 kills all in-flight requests immediately. The result is 500 errors during every rolling deployment — which happens multiple times per day in a busy organisation.

**`signal.NotifyContext`** (Go 1.16+) is the idiomatic API:
```go
ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
defer stop()
<-ctx.Done() // blocks until signal arrives
```

## References

- [signal.NotifyContext](https://pkg.go.dev/os/signal#NotifyContext)
- [http.Server.Shutdown](https://pkg.go.dev/net/http#Server.Shutdown)
- [Go Blog: Graceful Shutdown](https://go.dev/doc/articles/wiki/)
