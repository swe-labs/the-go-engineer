# Section 15: Time & Scheduling

## Beginner ? Expert Mapping

| Topic | Level | Importance | Engineering Concept |
| --- | --- | --- | --- |
| Timers | Beginner | High | Delayed execution |
| Tickers | Intermediate | Medium | Scheduled recurrence |
| Context | Advanced | Critical | Request-scoped deadlines |

## Engineering Depth

A common memory leak in Go is abandoning `time.Ticker` instances. The ticker uses a background
goroutine to push to its channel. If you do not call `ticker.Stop()`, that goroutine lives forever.
Contexts (`context.WithTimeout`) are the standard mechanism for passing scoped deadlines down call
stacks without relying on global tickers.

## References

1. [Package time](https://pkg.go.dev/time)

## Exercise: Console Reminder (`7-reminder`)

Build a countdown reminder that uses `time.NewTicker` and `time.AfterFunc`.

```bash
go run ./11-concurrency/time-and-scheduling/7-reminder/_starter 5 "Break time!"
go run ./11-concurrency/time-and-scheduling/7-reminder 5 "Break time!"
```

## Learning Path

| ID | Lesson | Concept | Requires |
| --- | --- | --- | --- |
| TM.1 | [time basics](./1-time) | `time.Time` � `Duration` � `Add` � `Sub` � wall vs monotonic clock | entry |
| TM.2 | [formatting](./2-formatting) | Reference Time `2006-01-02 15:04:05` � `Parse` � `RFC3339` | TM.1 |
| TM.3 | [timers & tickers](./3-timer-and-ticker) | `time.NewTimer` � `time.NewTicker` � `<-C` � `ticker.Stop()` leak | TM.1, TM.2 |
| TM.4 | [random numbers](./4-random) | `rand/v2` � `IntN` � `Shuffle` � `Perm` � seeded PCG | TM.1 |
| TM.5 | [scheduler](./5-schedule) | actor model � `ScheduleOnce` � `ScheduleInterval` � `StopAll` drain | TM.3 |
| TM.6 | [timezones](./6-timezone) | `time.LoadLocation` � IANA database � `In()` � always store UTC | TM.1, TM.2 |
