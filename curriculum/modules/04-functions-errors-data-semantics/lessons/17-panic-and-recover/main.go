package main

import (
	"errors"
	"fmt"
)

func safeDivide(a, b int) (result int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()
	if b == 0 {
		panic("division by zero")
	}
	return a / b, nil
}

func safeBatch(dividends, divisors []int) (results []int, err error) {
	if len(dividends) != len(divisors) {
		return nil, errors.New("mismatched slice lengths")
	}
	results = make([]int, len(dividends))
	for i := range dividends {
		func(i int) {
			res, e := safeDivide(dividends[i], divisors[i])
			if e != nil {
				err = errors.Join(err, fmt.Errorf("item %d: %w", i, e))
				return
			}
			results[i] = res
		}(i)
	}
	return results, err
}

func main() {
	testCases := []struct{ a, b int }{
		{10, 2},
		{10, 0},
		{20, 4},
	}
	for _, tc := range testCases {
		res, err := safeDivide(tc.a, tc.b)
		if err != nil {
			fmt.Printf("safeDivide(%d, %d): error: %v\n", tc.a, tc.b, err)
		} else {
			fmt.Printf("safeDivide(%d, %d): %d\n", tc.a, tc.b, res)
		}
	}

	// Batch test
	results, err := safeBatch([]int{10, 20, 30}, []int{2, 0, 3})
	if err != nil {
		fmt.Printf("safeBatch errors: %v\n", err)
	}
	fmt.Printf("safeBatch results: %v\n", results)
}
