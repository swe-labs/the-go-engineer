package main

import "testing"

// mockNotifier (Struct): groups the state used by the mock notifier example boundary.
type mockNotifier struct {
	called bool
}

// mockNotifier.Send (Method): applies the send operation to receiver state at a visible boundary.
func (m *mockNotifier) Send(_ string) { m.called = true }

// notifier (Interface): captures the behavior boundary the notifier example depends on.
type notifier interface{ Send(string) }

// notify (Function): runs the notify step and keeps its inputs, outputs, or errors visible.
func notify(n notifier, msg string) { n.Send(msg) }

func TestTE8Mock(t *testing.T) {
	mock := &mockNotifier{}
	notify(mock, "hello")
	if !mock.called {
		t.Fatal("expected notifier to be called")
	}
}
