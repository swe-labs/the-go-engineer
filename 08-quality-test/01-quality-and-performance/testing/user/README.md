# TE.1-TE.3 Testing Fundamentals

## Mission

Learn the foundational Go testing patterns through one small user-focused package.

This surface covers:

- `TE.1` unit testing
- `TE.2` table-driven tests
- `TE.3` HTTP handler testing

## Files

- [user.go](./user.go): code under test plus teaching comments
- [user_test.go](./user_test.go): basic tests and table-driven tests
- [greeting_test.go](./greeting_test.go): testable design with `io.Writer`
- [http_handler_test.go](./http_handler_test.go): handler testing with `httptest`

## Run Instructions

```bash
go test ./08-quality-test/01-quality-and-performance/testing/user
```

## Success Criteria

You should be able to:

- write small focused tests with `testing.T`
- structure table-driven tests with `t.Run`
- test handler behavior without starting a real server
- explain why dependency injection makes code easier to test

## Next Step

After this surface, continue to [TE.4 benchmarking](../benchmarks) or back to the
[Testing track](../README.md).
