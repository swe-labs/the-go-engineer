package main

import "testing"

func TestTE5Subtests(t *testing.T) {
	cases := map[string]string{
		"trim":  "  Go  ",
		"lower": "TEST",
	}

	for name, input := range cases {
		name, input := name, input
		t.Run(name, func(t *testing.T) {
			tmp := t.TempDir()
			t.Cleanup(func() {
				_ = tmp
			})
			if got := te_5Summary(input); got == "" {
				t.Fatalf("expected summary for %s", name)
			}
		})
	}
}
