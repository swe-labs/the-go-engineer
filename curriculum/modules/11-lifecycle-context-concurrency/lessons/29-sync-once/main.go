package main

import (
	"fmt"
	"sync"
)

type Config struct {
	Addr string
	Port int
}

var (
	config     *Config
	configOnce sync.Once
)

func LoadConfig() *Config {
	configOnce.Do(func() {
		fmt.Println("Loading configuration...")
		config = &Config{
			Addr: "0.0.0.0",
			Port: 8080,
		}
	})
	return config
}

type Pool struct {
	items chan interface{}
	once  sync.Once
}

func NewPool(size int) *Pool {
	p := &Pool{}
	p.once.Do(func() {
		p.items = make(chan interface{}, size)
		for i := 0; i < size; i++ {
			p.items <- struct{}{}
		}
	})
	return p
}

func (p *Pool) Get() interface{} {
	return <-p.items
}

func (p *Pool) Put(x interface{}) {
	p.items <- x
}

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cfg := LoadConfig()
			fmt.Printf("Config: %+v\n", cfg)
		}()
	}
	wg.Wait()

	pool := NewPool(3)
	fmt.Printf("Pool acquired %v\n", pool.Get())
	pool.Put(struct{}{})
	fmt.Println("Done")
}
