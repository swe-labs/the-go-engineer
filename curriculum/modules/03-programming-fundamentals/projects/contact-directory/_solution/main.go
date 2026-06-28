package main

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

type Contact struct {
	Name  string
	Email string
	Phone string
}

type Directory map[string]Contact

func NewContact(name, email, phone string) (Contact, error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)
	phone = strings.TrimSpace(phone)
	if name == "" || email == "" || phone == "" {
		return Contact{}, errors.New("all fields are required")
	}
	return Contact{Name: name, Email: email, Phone: phone}, nil
}

func AddContact(dir Directory, c Contact) error {
	if c.Phone == "" {
		return errors.New("phone is required")
	}
	if _, exists := dir[c.Phone]; exists {
		return fmt.Errorf("contact with phone %s already exists", c.Phone)
	}
	dir[c.Phone] = c
	return nil
}

func GetContact(dir Directory, phone string) (Contact, error) {
	c, exists := dir[phone]
	if !exists {
		return Contact{}, fmt.Errorf("contact with phone %s not found", phone)
	}
	return c, nil
}

func UpdateContact(dir Directory, phone string, c Contact) error {
	if _, exists := dir[phone]; !exists {
		return fmt.Errorf("contact with phone %s not found", phone)
	}
	dir[phone] = c
	return nil
}

func DeleteContact(dir Directory, phone string) error {
	if _, exists := dir[phone]; !exists {
		return fmt.Errorf("contact with phone %s not found", phone)
	}
	delete(dir, phone)
	return nil
}

func ListContacts(dir Directory) []Contact {
	if len(dir) == 0 {
		return nil
	}
	contacts := make([]Contact, 0, len(dir))
	for _, c := range dir {
		contacts = append(contacts, c)
	}
	sort.Slice(contacts, func(i, j int) bool {
		return strings.ToLower(contacts[i].Name) < strings.ToLower(contacts[j].Name)
	})
	return contacts
}

func SearchContacts(dir Directory, query string) []Contact {
	query = strings.ToLower(query)
	var results []Contact
	for _, c := range dir {
		if strings.Contains(strings.ToLower(c.Name), query) ||
			strings.Contains(strings.ToLower(c.Email), query) {
			results = append(results, c)
		}
	}
	if len(results) == 0 {
		return nil
	}
	sort.Slice(results, func(i, j int) bool {
		return strings.ToLower(results[i].Name) < strings.ToLower(results[j].Name)
	})
	return results
}

func main() {
	dir := make(Directory)

	alice, _ := NewContact("Alice", "alice@example.com", "555-0101")
	AddContact(dir, alice)
	fmt.Println("Added Alice")

	bob, _ := NewContact("Bob", "bob@example.com", "555-0102")
	AddContact(dir, bob)
	fmt.Println("Added Bob")

	charlie, _ := NewContact("Charlie", "charlie@example.com", "555-0103")
	AddContact(dir, charlie)
	fmt.Println("Added Charlie")

	contacts := ListContacts(dir)
	fmt.Println("--- All contacts ---")
	for _, c := range contacts {
		fmt.Printf("%s (%s) - %s\n", c.Name, c.Email, c.Phone)
	}

	fmt.Println("--- Search for 'alice' ---")
	for _, c := range SearchContacts(dir, "alice") {
		fmt.Printf("%s (%s) - %s\n", c.Name, c.Email, c.Phone)
	}

	DeleteContact(dir, "555-0102")
	fmt.Println("--- After deleting Bob ---")
	for _, c := range ListContacts(dir) {
		fmt.Printf("%s (%s) - %s\n", c.Name, c.Email, c.Phone)
	}
}
