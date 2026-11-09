package main

import (
	"fmt"
	"time"
)

func FetchData(source string) string {
	time.Sleep(10 * time.Millisecond)
	return "data from " + source
}

func Process(id int) (int, error) {
	if id <= 0 {
		return 0, fmt.Errorf("invalid id: %d", id)
	}
	return id * 2, nil
}

func main() {
	fmt.Println(FetchData("cache"))
	fmt.Println(Process(5))
}
