package main

import (
	"fmt"
)

// Contact represents a person in the directory.
type Contact struct {
	Name  string
	Email string
	Phone string
}

// Directory maps phone numbers to contacts.
type Directory map[string]Contact

// NewContact creates a Contact with trimmed fields.
// Returns an error if any field is empty after trimming.
func NewContact(name, email, phone string) (Contact, error) {
	// TODO: implement
	return Contact{}, nil
}

// AddContact adds a contact to the directory keyed by phone.
// Returns an error if the phone already exists or is empty.
func AddContact(dir Directory, c Contact) error {
	// TODO: implement
	return nil
}

// GetContact retrieves a contact by phone.
// Returns an error if not found.
func GetContact(dir Directory, phone string) (Contact, error) {
	// TODO: implement
	return Contact{}, nil
}

// UpdateContact replaces the contact at the given phone.
// Returns an error if the phone is not found.
func UpdateContact(dir Directory, phone string, c Contact) error {
	// TODO: implement
	return nil
}

// DeleteContact removes a contact by phone.
// Returns an error if the phone is not found.
func DeleteContact(dir Directory, phone string) error {
	// TODO: implement
	return nil
}

// ListContacts returns all contacts sorted by name (case-insensitive).
// Returns nil if empty.
func ListContacts(dir Directory) []Contact {
	// TODO: implement
	return nil
}

// SearchContacts returns contacts whose name or email contains
// the query string (case-insensitive). Returns nil if no matches.
func SearchContacts(dir Directory, query string) []Contact {
	// TODO: implement
	return nil
}

func main() {
	dir := make(Directory)

	alice, _ := NewContact("Alice", "alice@example.com", "555-0101")
	AddContact(dir, alice)
	fmt.Println("Added Alice")

	bob, _ := NewContact("Bob", "bob@example.com", "555-0102")
	AddContact(dir, bob)
	fmt.Println("Added Bob")

	contacts := ListContacts(dir)
	fmt.Println("--- All contacts ---")
	for _, c := range contacts {
		fmt.Printf("%s (%s) - %s\n", c.Name, c.Email, c.Phone)
	}
}
