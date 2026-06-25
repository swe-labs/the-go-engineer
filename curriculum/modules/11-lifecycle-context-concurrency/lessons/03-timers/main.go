package main

import (
	"fmt"
	"time"
)

func main() {
	timer := time.NewTimer(1 * time.Second)
	start := time.Now()
	<-timer.C
	fmt.Println("Timer fired after:", time.Since(start).Round(time.Millisecond))

	timer2 := time.NewTimer(500 * time.Millisecond)
	stopped := timer2.Stop()
	if stopped {
		fmt.Println("Timer stopped before firing")
	}

	timer3 := time.NewTimer(1 * time.Second)
	go func() {
		<-timer3.C
		fmt.Println("Timer 3 fired")
	}()
	time.Sleep(100 * time.Millisecond)
	reset := timer3.Reset(2 * time.Second)
	fmt.Println("Timer 3 reset, was active:", reset)
	time.Sleep(2500 * time.Millisecond)

	select {
	case <-time.After(100 * time.Millisecond):
		fmt.Println("time.After fired")
	}
}
