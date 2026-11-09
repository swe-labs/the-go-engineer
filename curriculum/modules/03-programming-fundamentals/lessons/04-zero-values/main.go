package main

import "fmt"

type Config struct {
	Host  string
	Port  int
	Debug bool
}

func main() {
	fmt.Println("=== basic zero values ===")
	var i int
	var f float64
	var s string
	var b bool
	fmt.Printf("int: %d, float64: %.1f, string: %q, bool: %t\n", i, f, s, b)

	fmt.Println("=== nil pointer ===")
	var p *int
	fmt.Printf("pointer: %v, p == nil: %t\n", p, p == nil)

	fmt.Println("=== nil slice ===")
	var nums []int
	fmt.Printf("nil slice: len=%d, cap=%d, isNil=%t\n", len(nums), cap(nums), nums == nil)
	for range nums {
	}
	nums = append(nums, 1)
	fmt.Println("after append:", nums)

	fmt.Println("=== nil map (read safe, write panics) ===")
	var m map[string]int
	fmt.Printf("nil map: len=%d, m == nil: %t\n", len(m), m == nil)
	_ = m["key"]

	fmt.Println("=== safeSet with nil map ===")
	m = safeSet(m, "alice", 95)
	m = safeSet(m, "bob", 87)
	val, ok := safeGet(m, "alice")
	fmt.Printf("alice: %d, found: %t\n", val, ok)
	val, ok = safeGet(m, "charlie")
	fmt.Printf("charlie: %d, found: %t\n", val, ok)

	fmt.Println("=== nil channel ===")
	var ch chan int
	fmt.Printf("nil channel: %v, ch == nil: %t\n", ch, ch == nil)

	fmt.Println("=== zero value struct ===")
	var cfg Config
	fmt.Printf("zero struct: Host=%q, Port=%d, Debug=%t\n", cfg.Host, cfg.Port, cfg.Debug)

	fmt.Println("=== safe range over nil slice ===")
	var items []string
	count := 0
	for range items {
		count++
	}
	fmt.Println("iterations over nil slice:", count)
	items = append(items, "a", "b")
	fmt.Println("after append:", items)
}

func safeSet(m map[string]int, key string, value int) map[string]int {
	if m == nil {
		m = make(map[string]int)
	}
	m[key] = value
	return m
}

func safeGet(m map[string]int, key string) (int, bool) {
	if m == nil {
		return 0, false
	}
	val, ok := m[key]
	return val, ok
}
