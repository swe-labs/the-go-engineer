# Section 1: Language Basics

## Beginner → Expert Mapping

| Topic | Level | Importance | Engineering Concept |
|-------|-------|------------|---------------------|
| Variables & Types | Beginner | High | Static typing, Zero values |
| Constants | Beginner | Low | Immutable values |
| Constants & iota | Intermediate | Medium | Enum simulation, bitmasking |

## Engineering Depth
Go is a strictly typed language. Unlike dynamic languages, Go defines memory bounds at compile time. One of Go's unique features is the **Zero Value**. If you declare `var x int`, Go automatically zeroes the memory, meaning `x` is predictably `0`, not undefined. This removes a whole class of uninitialized variable bugs.

## References
1. **[Tour of Go]** [Variables](https://go.dev/tour/basics/8)
2. **[Effective Go]** [Constants](https://go.dev/doc/effective_go#constants)

---

## 🏗 Exercise: Application Logger (`4-application-logger`)

This introductory project proves your understanding of types and formats.

### Step-by-Step Instructions & Hints
1. **Define your constants:** Create severities (`INFO`, `WARNING`, `ERROR`) using the `iota` enumerator.
2. **Create a log struct:** Build an `AppLog` type with fields for timestamp, severity, and message.
3. **Write the formatter:** Create a `PrintLog()` function that uses `fmt.Printf` to format your log line cleanly.
   - *Hint:* Use `%s` for strings, `%d` for numbers, and `%-10s` to right-pad text for perfect alignment!
