# Section 15: Time & Scheduling

## Beginner → Expert Mapping

| Topic | Level | Importance | Engineering Concept |
|-------|-------|------------|---------------------|
| Timers | Beginner | High | Delayed execution |
| Tickers | Intermediate| Medium | Scheduled recurrence |
| Context | Advanced | **Critical** | Request-scoped deadlines |

## Engineering Depth
A common memory leak in Go is abandoning `time.Ticker` instances. The ticker uses a background goroutine to push to its channel. If you do not call `ticker.Stop()`, that goroutine lives forever. Contexts (`context.WithTimeout`) are the Google standard mechanism for passing scoped deadlines down massive call stacks without relying on global tickers.

## References
1. **[Go Docs]** [Package time](https://pkg.go.dev/time)

---

## 🏗 Exercise: Console Reminder (`7-reminder`)

Build a countdown reminder that uses `time.NewTicker` and `time.AfterFunc`. Try it yourself first!

```bash
go run ./15-time-and-scheduling/7-reminder/_starter 5 "Break time!"  # Try the exercise
go run ./15-time-and-scheduling/7-reminder 5 "Break time!"           # See the solution
```
