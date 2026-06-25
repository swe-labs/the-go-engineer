package main

import (
	"encoding/json"
	"errors"
	"fmt"
)

type User struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email,omitempty"`
}

func MarshalUser(name string, age int, email string) (string, error) {
	if age < 0 {
		return "", errors.New("age must not be negative")
	}
	u := User{Name: name, Age: age, Email: email}
	b, err := json.MarshalIndent(u, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func main() {
	cases := []struct {
		name  string
		age   int
		email string
	}{
		{"Alice", 30, "alice@example.com"},
		{"Bob", 25, ""},
		{"Carol", -1, "carol@test.org"},
	}
	for _, c := range cases {
		result, err := MarshalUser(c.name, c.age, c.email)
		if err != nil {
			fmt.Printf("MarshalUser(%q, %d, %q) error: %v\n", c.name, c.age, c.email, err)
		} else {
			fmt.Printf("MarshalUser(%q, %d, %q) =\n%s\n", c.name, c.age, c.email, result)
		}
	}
}
