# Go Memory Model & Allocation

The **Go Memory Model** specifies the conditions under which reads of a variable in one goroutine can be guaranteed to observe values produced by writes to the same variable in another goroutine.

## Stack vs. Heap Allocation

Go uses **escape analysis** performed at compile-time to determine whether a variable can be safely allocated on the stack or must escape to the heap.

```go
func createOnStack() int {
	x := 42
	return x // x does not escape, kept on stack
}

func escapeToHeap() *int {
	y := 100
	return &y // y escapes because pointer outlives function frame
}
```

## Happens-Before Guarantees

- **Package Initialization**: If package `p` imports package `q`, the completion of `q`'s `init` functions happens before the start of any of `p`'s.
- **Goroutine Creation**: The `go` statement that starts a new goroutine happens before the execution of the goroutine begins.
- **Channel Operations**: A send on a channel happens before the corresponding receive from that channel completes.
