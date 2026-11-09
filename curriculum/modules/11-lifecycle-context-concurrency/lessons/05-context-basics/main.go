package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx := context.Background()
	fmt.Println("Background ctx err:", ctx.Err())
	fmt.Println("Background ctx done:", ctx.Done())

	ctxCancel, cancel := context.WithCancel(ctx)
	cancel()
	fmt.Println("Canceled ctx err:", ctxCancel.Err())

	ctxTimeout, timeoutCancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer timeoutCancel()
	<-ctxTimeout.Done()
	fmt.Println("Timeout ctx err:", ctxTimeout.Err())
}
