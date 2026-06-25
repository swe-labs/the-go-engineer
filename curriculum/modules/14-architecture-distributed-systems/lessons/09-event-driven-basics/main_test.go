package main

import (
	"sync"
	"testing"
)

func TestEventBus_PublishSubscribe(t *testing.T) {
	bus := NewEventBus()
	var got []string
	var mu sync.Mutex

	bus.Subscribe("test.event", func(e Event) {
		mu.Lock()
		got = append(got, e.Data.(string))
		mu.Unlock()
	})

	bus.Publish(Event{Type: "test.event", Data: "hello"})
	bus.Publish(Event{Type: "test.event", Data: "world"})

	if len(got) != 2 {
		t.Errorf("expected 2 events, got %d", len(got))
	}
}

func TestEventBus_MultipleHandlers(t *testing.T) {
	bus := NewEventBus()
	var count1, count2 int
	var mu sync.Mutex

	bus.Subscribe("evt", func(e Event) {
		mu.Lock()
		count1++
		mu.Unlock()
	})
	bus.Subscribe("evt", func(e Event) {
		mu.Lock()
		count2++
		mu.Unlock()
	})

	bus.Publish(Event{Type: "evt", Data: nil})

	if count1 != 1 {
		t.Errorf("handler 1 called %d times, expected 1", count1)
	}
	if count2 != 1 {
		t.Errorf("handler 2 called %d times, expected 1", count2)
	}
}

func TestEventBus_UnrelatedEvent(t *testing.T) {
	bus := NewEventBus()
	called := false
	bus.Subscribe("type.a", func(e Event) {
		called = true
	})
	bus.Publish(Event{Type: "type.b", Data: nil})
	if called {
		t.Error("handler should not be called for unrelated event type")
	}
}

func TestNotificationService_OnOrderPlaced(t *testing.T) {
	bus := NewEventBus()
	ns := &NotificationService{}
	bus.Subscribe("order.placed", ns.OnOrderPlaced)

	bus.Publish(Event{Type: "order.placed", Data: OrderPlacedData{OrderID: "ORD-001", Total: 49.99}})

	notifs := ns.Notifications()
	if len(notifs) != 1 {
		t.Fatalf("expected 1 notification, got %d", len(notifs))
	}
}

func TestInventoryService_OnOrderPlaced(t *testing.T) {
	bus := NewEventBus()
	is := &InventoryService{}
	bus.Subscribe("order.placed", is.OnOrderPlaced)

	bus.Publish(Event{Type: "order.placed", Data: OrderPlacedData{OrderID: "ORD-001", Total: 49.99}})

	reserved := is.Reserved()
	if len(reserved) != 1 {
		t.Fatalf("expected 1 reservation, got %d", len(reserved))
	}
}
