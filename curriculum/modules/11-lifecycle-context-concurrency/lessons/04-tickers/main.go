package main

import (
	"fmt"
	"time"
)

func main() {
	ticker := time.NewTicker(200 * time.Millisecond)
	done := make(chan bool)

	go func() {
		for i := 0; i < 5; i++ {
			fmt.Println("Tick at:", time.Now().Format("15:04:05.000"))
			<-ticker.C
		}
		ticker.Stop()
		done <- true
	}()

	<-done
	fmt.Println("Ticker stopped")
}
