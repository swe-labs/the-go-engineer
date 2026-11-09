package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func LoadUsers() ([]User, error) {
	data, err := os.ReadFile(filepath.Join("testdata", "users.json"))
	if err != nil {
		return nil, err
	}
	var users []User
	if err := json.Unmarshal(data, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func main() {
	users, err := LoadUsers()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for _, u := range users {
		fmt.Printf("%s is %d years old\n", u.Name, u.Age)
	}
}
