package main

import (
	"fmt"
	"sync"
)

type Event struct {
	Type string
	Data interface{}
}

type EventHandler func(Event)

type EventBus struct {
	mu       sync.RWMutex
	handlers map[string][]EventHandler
}

func NewEventBus() *EventBus {
	return &EventBus{handlers: make(map[string][]EventHandler)}
}

func (b *EventBus) Subscribe(eventType string, handler EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventType] = append(b.handlers[eventType], handler)
}

func (b *EventBus) Publish(event Event) {
	b.mu.RLock()
	handlers := b.handlers[event.Type]
	b.mu.RUnlock()
	for _, h := range handlers {
		h(event)
	}
}

type OrderPlacedData struct {
	OrderID string
	Total   float64
}

type NotificationService struct {
	notifications []string
	mu            sync.Mutex
}

func (n *NotificationService) OnOrderPlaced(event Event) {
	data := event.Data.(OrderPlacedData)
	n.mu.Lock()
	n.notifications = append(n.notifications, fmt.Sprintf("Order %s for $%.2f placed", data.OrderID, data.Total))
	n.mu.Unlock()
}

func (n *NotificationService) Notifications() []string {
	n.mu.Lock()
	defer n.mu.Unlock()
	result := make([]string, len(n.notifications))
	copy(result, n.notifications)
	return result
}

type InventoryService struct {
	reserved []string
	mu       sync.Mutex
}

func (i *InventoryService) OnOrderPlaced(event Event) {
	data := event.Data.(OrderPlacedData)
	i.mu.Lock()
	i.reserved = append(i.reserved, fmt.Sprintf("Inventory reserved for %s", data.OrderID))
	i.mu.Unlock()
}

func (i *InventoryService) Reserved() []string {
	i.mu.Lock()
	defer i.mu.Unlock()
	result := make([]string, len(i.reserved))
	copy(result, i.reserved)
	return result
}

func main() {
	bus := NewEventBus()
	notifier := &NotificationService{}
	inventory := &InventoryService{}

	bus.Subscribe("order.placed", notifier.OnOrderPlaced)
	bus.Subscribe("order.placed", inventory.OnOrderPlaced)

	bus.Publish(Event{
		Type: "order.placed",
		Data: OrderPlacedData{OrderID: "ORD-001", Total: 99.99},
	})

	fmt.Println("Notifications:", notifier.Notifications())
	fmt.Println("Inventory:", inventory.Reserved())
}
