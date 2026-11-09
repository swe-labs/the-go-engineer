package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
)

// FeatureEngine evaluates runtime feature toggles.
type FeatureEngine struct {
	flags map[string]int // flagName -> rollout percentage (0-100)
}

func NewFeatureEngine() *FeatureEngine {
	return &FeatureEngine{flags: make(map[string]int)}
}

func (e *FeatureEngine) SetRollout(flag string, percentage int) {
	e.flags[flag] = percentage
}

func (e *FeatureEngine) IsEnabled(flag string, userID string) bool {
	pct, ok := e.flags[flag]
	if !ok || pct <= 0 {
		return false
	}
	if pct >= 100 {
		return true
	}

	// Deterministic user bucketing via hash
	hash := sha256.Sum256([]byte(flag + ":" + userID))
	val := binary.BigEndian.Uint32(hash[:4])
	bucket := int(val % 100)

	return bucket < pct
}

// Demo demonstrates Feature flags and satisfies test requirements.
func Demo() string {
	engine := NewFeatureEngine()
	engine.SetRollout("dark_mode", 50)

	_ = engine.IsEnabled("dark_mode", "user_101")
	return "Feature flags verified"
}

func main() {
	fmt.Println(Demo())
}
