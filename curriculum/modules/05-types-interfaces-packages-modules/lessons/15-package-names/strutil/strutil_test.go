package strutil

import "testing"

func TestUpper(t *testing.T) {
	if got := Upper("hello"); got != "HELLO" {
		t.Errorf("Upper(hello) = %q; want HELLO", got)
	}
	if got := Upper("Go"); got != "GO" {
		t.Errorf("Upper(Go) = %q; want GO", got)
	}
	if got := Upper(""); got != "" {
		t.Errorf("Upper(empty) = %q; want empty", got)
	}
}
