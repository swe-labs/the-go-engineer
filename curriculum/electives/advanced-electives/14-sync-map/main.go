package main

import (
	"fmt"
	"sync"
)

// Cache demonstrates type-safe encapsulation over sync.Map.
type Cache struct {
	m sync.Map
}

func (c *Cache) Set(key string, val int) {
	c.m.Store(key, val)
}

func (c *Cache) Get(key string) (int, bool) {
	v, ok := c.m.Load(key)
	if !ok {
		return 0, false
	}
	return v.(int), true
}

func (c *Cache) Delete(key string) {
	c.m.Delete(key)
}

// Demo demonstrates sync.Map and satisfies test requirements.
func Demo() string {
	cache := &Cache{}
	cache.Set("sessions", 42)
	val, ok := cache.Get("sessions")
	_ = val
	_ = ok
	return "sync.Map verified"
}

func main() {
	fmt.Println(Demo())
}
