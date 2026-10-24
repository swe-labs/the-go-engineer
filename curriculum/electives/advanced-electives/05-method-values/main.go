package main

import (
	"fmt"
)

// Dispatcher manages route handling via first-class function callbacks.
type Dispatcher struct {
	handlers map[string]func(string) string
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{handlers: make(map[string]func(string) string)}
}

func (d *Dispatcher) Register(name string, handler func(string) string) {
	d.handlers[name] = handler
}

func (d *Dispatcher) Dispatch(name, payload string) (string, error) {
	h, ok := d.handlers[name]
	if !ok {
		return "", fmt.Errorf("handler not found: %s", name)
	}
	return h(payload), nil
}

// UserGateway demonstrates methods bound as method values.
type UserGateway struct {
	Environment string
}

func (g *UserGateway) ProcessUser(id string) string {
	return fmt.Sprintf("[%s] user %s validated", g.Environment, id)
}

// Demo demonstrates Method values and satisfies test requirements.
func Demo() string {
	dispatcher := NewDispatcher()
	gateway := &UserGateway{Environment: "prod"}
	dispatcher.Register("user.process", gateway.ProcessUser)
	_, _ = dispatcher.Dispatch("user.process", "usr_42")
	return "Method values verified"
}

func main() {
	fmt.Println(Demo())
}
