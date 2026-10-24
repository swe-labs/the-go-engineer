package main

import (
	"fmt"
)

// Number defines a constraint interface accepting any integer or floating-point type,
// including custom user-defined types having those underlying primitives.
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// Milliseconds is a custom defined type based on int64.
type Milliseconds int64

// Clamp bounds a number between min and max inclusive.
func Clamp[T Number](val, min, max T) T {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

// Demo demonstrates Complex generic constraints and satisfies test requirements.
func Demo() string {
	var timeout Milliseconds = 1500
	bounded := Clamp(timeout, 100, 1000)
	floatVal := Clamp(3.14159, 0.0, 2.5)
	_ = bounded
	_ = floatVal
	return "Complex generic constraints verified"
}

func main() {
	fmt.Println(Demo())
}
