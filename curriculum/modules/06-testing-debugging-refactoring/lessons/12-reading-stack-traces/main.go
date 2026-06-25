package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Println("=== Walk the stack with runtime.Caller ===")
	First()

	fmt.Println("\n=== Capture full trace with runtime.Stack ===")
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	fmt.Printf("%s\n", buf[:n])

	fmt.Println("\n=== Panic and recover ===")
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered:", r)
			}
		}()
		trigger()
	}()
}

func First()  { Second() }
func Second() { Third() }
func Third() {
	for skip := 0; ; skip++ {
		pc, file, line, ok := runtime.Caller(skip)
		if !ok {
			break
		}
		fn := runtime.FuncForPC(pc)
		fmt.Printf("  skip=%d  %s (%s:%d)\n", skip, fn.Name(), file, line)
	}
}

func trigger() { deep() }
func deep()    { panic("boom") }

func GetCallerInfo(skip int) (funcName, file string, line int, ok bool) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "", "", 0, false
	}
	fn := runtime.FuncForPC(pc)
	return fn.Name(), file, line, true
}

func StackDepth() int {
	for i := 0; ; i++ {
		_, _, _, ok := runtime.Caller(i)
		if !ok {
			return i
		}
	}
}
