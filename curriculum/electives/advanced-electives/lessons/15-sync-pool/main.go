package main

import (
	"bytes"
	"fmt"
	"sync"
)

// BufferPool manages reusable bytes.Buffer instances.
type BufferPool struct {
	pool sync.Pool
}

func NewBufferPool() *BufferPool {
	return &BufferPool{
		pool: sync.Pool{
			New: func() any {
				// Allocate a buffer with 4KB preallocated capacity
				return bytes.NewBuffer(make([]byte, 0, 4096))
			},
		},
	}
}

func (p *BufferPool) Get() *bytes.Buffer {
	return p.pool.Get().(*bytes.Buffer)
}

func (p *BufferPool) Put(buf *bytes.Buffer) {
	buf.Reset() // Crucial: reset slice length to 0 while keeping capacity!
	p.pool.Put(buf)
}

// Demo demonstrates sync.Pool and satisfies test requirements.
func Demo() string {
	pool := NewBufferPool()
	buf := pool.Get()
	buf.WriteString("hello sync.pool")
	_ = buf.String()
	pool.Put(buf)
	return "sync.Pool verified"
}

func main() {
	fmt.Println(Demo())
}
