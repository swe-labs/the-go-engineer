package main

import "testing"

func FuzzTE6Summary(f *testing.F) {
	for _, seed := range []string{"Go", "  API  ", "TEST"} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		got := te_6Summary(input)
		if got != te_6Summary(got) {
			t.Fatalf("summary should be idempotent")
		}
	})
}
