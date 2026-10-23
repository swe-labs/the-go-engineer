package main

import (
	"cmp"
	"fmt"
	"slices"
	"time"
)

// Transaction represents a financial or business ledger entry.
type Transaction struct {
	ID        string
	Priority  int
	Amount    float64
	Timestamp time.Time
}

// DeduplicatePreserveOrder removes duplicate strings while maintaining first-seen order.
// Complexity: O(N) time, O(N) space.
func DeduplicatePreserveOrder(items []string) []string {
	seen := make(map[string]struct{}, len(items))
	result := make([]string, 0, len(items))
	for _, item := range items {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

// PartitionActive in-place partitions items so matching items appear first,
// returning the count of matching items. Demonstrates two-pointer technique.
// Complexity: O(N) time, O(1) auxiliary space.
func PartitionActive(items []int, isTarget func(int) bool) int {
	left := 0
	for right := 0; right < len(items); right++ {
		if isTarget(items[right]) {
			items[left], items[right] = items[right], items[left]
			left++
		}
	}
	return left
}

// FindPriceTier finds the price for a quantity using binary search on sorted quantity thresholds.
// Complexity: O(log N) time, O(1) space.
func FindPriceTier(thresholds []int, prices []float64, qty int) (float64, bool) {
	if len(thresholds) == 0 || qty < thresholds[0] {
		return 0, false
	}
	idx, exact := slices.BinarySearch(thresholds, qty)
	if exact {
		return prices[idx], true
	}
	// If not exact, idx is insertion point; previous threshold applies
	return prices[idx-1], true
}

// SortTransactions sorts transactions by Priority (descending) then by Timestamp (ascending).
// Complexity: O(N log N) time via pdqsort.
func SortTransactions(txs []Transaction) {
	slices.SortFunc(txs, func(a, b Transaction) int {
		// Priority descending (higher priority first)
		if c := cmp.Compare(b.Priority, a.Priority); c != 0 {
			return c
		}
		// Timestamp ascending (older first)
		return a.Timestamp.Compare(b.Timestamp)
	})
}

func main() {
	fmt.Println("=== 1. Order-Preserving Deduplication ===")
	raw := []string{"apple", "banana", "apple", "cherry", "banana", "date"}
	deduped := DeduplicatePreserveOrder(raw)
	fmt.Printf("Original: %v\nDeduped:  %v\n", raw, deduped)

	fmt.Println("\n=== 2. Two-Pointer Partitioning ===")
	numbers := []int{3, 8, 2, 5, 12, 7, 4, 1}
	count := PartitionActive(numbers, func(n int) bool { return n%2 == 0 })
	fmt.Printf("Partitioned (evens first): %v (even count = %d)\n", numbers, count)

	fmt.Println("\n=== 3. Binary Search Price Tiers ===")
	thresholds := []int{1, 10, 50, 100, 500}
	prices := []float64{15.0, 12.5, 10.0, 8.0, 6.0}
	for _, q := range []int{1, 5, 25, 150, 600} {
		price, found := FindPriceTier(thresholds, prices, q)
		if found {
			fmt.Printf("Qty: %3d -> Unit Price: $%.2f\n", q, price)
		}
	}

	fmt.Println("\n=== 4. Custom Multi-Field Sorting (slices.SortFunc) ===")
	now := time.Now()
	txs := []Transaction{
		{ID: "tx-1", Priority: 1, Amount: 100.0, Timestamp: now.Add(2 * time.Minute)},
		{ID: "tx-2", Priority: 3, Amount: 250.0, Timestamp: now.Add(1 * time.Minute)},
		{ID: "tx-3", Priority: 1, Amount: 50.0, Timestamp: now.Add(1 * time.Minute)},
		{ID: "tx-4", Priority: 2, Amount: 500.0, Timestamp: now},
	}
	SortTransactions(txs)
	for _, tx := range txs {
		fmt.Printf("ID: %-5s | Priority: %d | Time: %s\n", tx.ID, tx.Priority, tx.Timestamp.Format("15:04:05"))
	}
}
