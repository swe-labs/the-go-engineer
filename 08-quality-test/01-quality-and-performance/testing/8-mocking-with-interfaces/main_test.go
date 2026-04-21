package main

import "testing"

type mockNotifier struct {
	called bool
}

func (m *mockNotifier) Send(_ string) { m.called = true }

type notifier interface{ Send(string) }

func notify(n notifier, msg string) { n.Send(msg) }

func TestTE8Mock(t *testing.T) {
	mock := &mockNotifier{}
	notify(mock, "hello")
	if !mock.called {
		t.Fatal("expected notifier to be called")
	}
}
