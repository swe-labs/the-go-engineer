package main

import (
	"errors"
	"fmt"
)

var ErrTimeout = errors.New("timeout")

type NetworkError struct {
	Msg       string
	IsTimeout bool
}

func (n *NetworkError) Error() string {
	return n.Msg
}

func (n *NetworkError) Is(target error) bool {
	return target == ErrTimeout && n.IsTimeout
}

func callService(shouldTimeout bool) error {
	if shouldTimeout {
		return &NetworkError{Msg: "service timed out", IsTimeout: true}
	}
	return &NetworkError{Msg: "service failed", IsTimeout: false}
}

func main() {
	for _, val := range []bool{true, false} {
		err := callService(val)
		if errors.Is(err, ErrTimeout) {
			fmt.Printf("callService(%v): timeout error: %v\n", val, err)
		} else {
			fmt.Printf("callService(%v): other error: %v\n", val, err)
		}
	}
}
