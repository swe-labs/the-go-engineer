package main

import "fmt"

func splitAndCopy(data []int, chunkSize int) [][]int {
	var result [][]int
	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}
		chunk := make([]int, end-i)
		copy(chunk, data[i:end])
		result = append(result, chunk)
	}
	return result
}

func unsafeSplit(data []int, chunkSize int) [][]int {
	var result [][]int
	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}
		result = append(result, data[i:end])
	}
	return result
}

func main() {
	data := []int{0, 1, 2, 3, 4, 5}

	safe := splitAndCopy(data, 2)
	unsafe := unsafeSplit(data, 2)

	data[1] = 99
	data[3] = 88
	data[5] = 77

	fmt.Println("original modified:", data)
	fmt.Println("safe chunks (no aliasing):", safe)
	fmt.Println("unsafe chunks (aliasing):", unsafe)
}
