package main

import "fmt"

// UserStore stores user information.
type UserStore interface {
	Save(name string, age int) error
	Find(name string) (int, error)
}

// RealUserStore is a production implementation backed by a map.
type RealUserStore struct {
	db map[string]int
}

func (r *RealUserStore) Save(name string, age int) error {
	r.db[name] = age
	return nil
}

func (r *RealUserStore) Find(name string) (int, error) {
	age, ok := r.db[name]
	if !ok {
		return 0, fmt.Errorf("user %q not found", name)
	}
	return age, nil
}

// Calculator performs arithmetic operations.
type Calculator interface {
	Add(a, b int) int
	Multiply(a, b int) int
}

// RealCalculator performs real arithmetic.
type RealCalculator struct{}

func (RealCalculator) Add(a, b int) int      { return a + b }
func (RealCalculator) Multiply(a, b int) int { return a * b }

// Double returns n*2 using Add.
func Double(c Calculator, n int) int {
	return c.Add(n, n)
}

// Greeter greets users from a store.
type Greeter struct {
	store UserStore
}

func NewGreeter(store UserStore) *Greeter {
	return &Greeter{store: store}
}

func (g *Greeter) Greet(name string) (string, error) {
	age, err := g.store.Find(name)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Hello, %s! You are %d years old.", name, age), nil
}

func main() {
	store := &RealUserStore{db: make(map[string]int)}
	store.Save("Alice", 30)
	g := NewGreeter(store)
	msg, _ := g.Greet("Alice")
	fmt.Println(msg)

	fmt.Println("Double(5) =", Double(RealCalculator{}, 5))
}
