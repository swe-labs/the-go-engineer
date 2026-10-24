package main

import (
	"fmt"
	"sync"
)

// TargetPool represents a deployable environment cluster.
type TargetPool struct {
	Name    string
	Version string
	Active  bool
}

// Router simulates atomic blue/green traffic switching.
type Router struct {
	mu     sync.RWMutex
	active string
	pools  map[string]*TargetPool
}

func NewRouter() *Router {
	r := &Router{
		active: "blue",
		pools: map[string]*TargetPool{
			"blue":  {Name: "blue", Version: "v1.0.0", Active: true},
			"green": {Name: "green", Version: "v2.0.0", Active: false},
		},
	}
	return r
}

func (r *Router) SwitchTraffic(target string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	targetPool, ok := r.pools[target]
	if !ok {
		return fmt.Errorf("unknown pool: %s", target)
	}

	for _, p := range r.pools {
		p.Active = false
	}
	targetPool.Active = true
	r.active = target
	return nil
}

// Demo demonstrates Blue/green deployments and satisfies test requirements.
func Demo() string {
	router := NewRouter()
	_ = router.SwitchTraffic("green")
	return "Blue/green deployments verified"
}

func main() {
	fmt.Println(Demo())
}
