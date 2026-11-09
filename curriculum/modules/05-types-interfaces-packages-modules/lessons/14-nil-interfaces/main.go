package main

import "fmt"

type MyError struct {
	Msg string
}

func (e *MyError) Error() string {
	if e == nil {
		return "<nil>"
	}
	return e.Msg
}

func doWork() error {
	var err *MyError
	return err
}

func doWorkSafe() error {
	return nil
}

func getError(shouldFail bool) error {
	if shouldFail {
		return &MyError{Msg: "something went wrong"}
	}
	return nil
}

func main() {
	err := doWork()
	fmt.Printf("doWork: err == nil = %v, err = %v\n", err == nil, err)

	err2 := doWorkSafe()
	fmt.Printf("doWorkSafe: err == nil = %v, err = %v\n", err2 == nil, err2)

	fmt.Println("--- getError ---")
	for _, flag := range []bool{false, true} {
		e := getError(flag)
		fmt.Printf("getError(%v) == nil: %v, value: %v\n", flag, e == nil, e)
	}
}
