package main

import (
	"flag"
	"fmt"
)

func main() {
	name := flag.String("name", "World", "a name to greet")
	count := flag.Int("count", 1, "number of times to greet")
	verbose := flag.Bool("verbose", false, "enable verbose output")

	flag.Parse()

	if *verbose {
		fmt.Printf("(verbose) name=%s count=%d\n", *name, *count)
	}

	for i := 0; i < *count; i++ {
		fmt.Printf("Hello, %s!\n", *name)
	}

	fmt.Println("Non-flag args:", flag.Args())
}
