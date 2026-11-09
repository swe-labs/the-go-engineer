package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now()
	fmt.Println("Now (Unix):", now.Unix())
	fmt.Println("Now (RFC3339):", now.Format(time.RFC3339))
	fmt.Println("Now (custom):", now.Format("2006-01-02 15:04:05"))

	layout := "2006-01-02"
	parsed, err := time.Parse(layout, "2026-06-01")
	if err != nil {
		panic(err)
	}
	fmt.Println("Parsed:", parsed.Format(time.RFC3339))

	then := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	fmt.Println("Since then:", time.Since(then))
	fmt.Println("Until next year:", time.Until(then.AddDate(1, 0, 0)))
}
