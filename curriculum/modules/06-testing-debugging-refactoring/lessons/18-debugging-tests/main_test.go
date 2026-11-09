package main

import "testing"

func TestProcessOrder(t *testing.T) {
	tests := []struct {
		name     string
		items    []string
		wantErr  bool
		wantCost int
	}{
		{"valid order", []string{"item1", "item2"}, false, 200},
		{"empty order", []string{}, true, 0},
		{"unknown item", []string{"nope"}, true, 0},
		{"single item", []string{"item1"}, false, 100},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cost, err := processOrder(tc.items)
			if (err != nil) != tc.wantErr {
				t.Errorf("processOrder() error = %v, wantErr %v", err, tc.wantErr)
			}
			if cost != tc.wantCost {
				t.Errorf("processOrder() cost = %d, want %d", cost, tc.wantCost)
			}
		})
	}
}
