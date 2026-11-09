package main

import "testing"

func TestClosingChannelsCompiles(t *testing.T) {
}

func TestReceiveFromClosedChannel(t *testing.T) {
	ch := make(chan int, 1)
	ch <- 42
	close(ch)
	v, ok := <-ch
	if v != 42 || !ok {
		t.Fatalf("expected 42, true; got %d, %t", v, ok)
	}
	v, ok = <-ch
	if v != 0 || ok {
		t.Fatalf("expected 0, false; got %d, %t", v, ok)
	}
}

func TestRangeOverClosedChannel(t *testing.T) {
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3
	close(ch)

	var got []int
	for v := range ch {
		got = append(got, v)
	}
	if len(got) != 3 {
		t.Fatalf("expected 3 values, got %d", len(got))
	}
}

func TestCanCloseNilChannel(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on closing nil channel")
		}
	}()
	var ch chan int
	close(ch)
}

func TestPanicOnDoubleClose(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on double close")
		}
	}()
	ch := make(chan int)
	close(ch)
	close(ch)
}
