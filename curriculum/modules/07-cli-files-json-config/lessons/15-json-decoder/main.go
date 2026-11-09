package main

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

func SumJSON(r io.Reader) (int, error) {
	dec := json.NewDecoder(r)
	var sum int
	for {
		var v int
		err := dec.Decode(&v)
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}
		sum += v
	}
	return sum, nil
}

func main() {
	data := `42
100
8
`
	sum, err := SumJSON(strings.NewReader(data))
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Sum: %d\n", sum)
	}

	// Error case: non-number
	bad := `42
not-a-number
`
	_, err = SumJSON(strings.NewReader(bad))
	fmt.Println("Bad input error:", err)
}
