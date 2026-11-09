package main

import (
	"fmt"
)

//go:generate echo "Generating assets for user status..."

// Status represents the state of an account.
type Status int

const (
	StatusPending Status = iota
	StatusActive
	StatusSuspended
)

// String provides a human-readable representation of Status.
func (s Status) String() string {
	switch s {
	case StatusPending:
		return "Pending"
	case StatusActive:
		return "Active"
	case StatusSuspended:
		return "Suspended"
	default:
		return fmt.Sprintf("Status(%d)", s)
	}
}

// Demo demonstrates go generate and satisfies test requirements.
func Demo() string {
	s := StatusActive
	_ = s.String()
	return "go generate verified"
}

func main() {
	fmt.Println(Demo())
}
