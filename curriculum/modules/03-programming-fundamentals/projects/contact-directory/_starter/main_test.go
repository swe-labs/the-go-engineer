package main

import "testing"

func TestNewContact(t *testing.T) {
	c, err := NewContact("Alice", "alice@example.com", "555-0101")
	if err != nil {
		t.Fatalf("NewContact unexpected error: %v", err)
	}
	if c.Name != "Alice" || c.Email != "alice@example.com" || c.Phone != "555-0101" {
		t.Errorf("NewContact = %+v; want {Alice alice@example.com 555-0101}", c)
	}
}

func TestNewContactEmptyField(t *testing.T) {
	_, err := NewContact("", "alice@example.com", "555-0101")
	if err == nil {
		t.Error("NewContact with empty name expected error")
	}
	_, err = NewContact("Alice", "", "555-0101")
	if err == nil {
		t.Error("NewContact with empty email expected error")
	}
	_, err = NewContact("Alice", "alice@example.com", "")
	if err == nil {
		t.Error("NewContact with empty phone expected error")
	}
}

func TestNewContactTrim(t *testing.T) {
	c, err := NewContact("  Alice  ", "  alice@example.com  ", "  555-0101  ")
	if err != nil {
		t.Fatalf("NewContact trim unexpected error: %v", err)
	}
	if c.Name != "Alice" {
		t.Errorf("NewContact trimmed name = %q; want %q", c.Name, "Alice")
	}
}

func TestAddAndGetContact(t *testing.T) {
	dir := make(Directory)
	c, _ := NewContact("Alice", "alice@example.com", "555-0101")
	err := AddContact(dir, c)
	if err != nil {
		t.Fatalf("AddContact unexpected error: %v", err)
	}
	got, err := GetContact(dir, "555-0101")
	if err != nil {
		t.Fatalf("GetContact unexpected error: %v", err)
	}
	if got.Name != "Alice" {
		t.Errorf("GetContact name = %q; want %q", got.Name, "Alice")
	}
}

func TestAddContactDuplicate(t *testing.T) {
	dir := make(Directory)
	c1, _ := NewContact("Alice", "a@x.com", "555-0101")
	c2, _ := NewContact("Bob", "b@x.com", "555-0101")
	AddContact(dir, c1)
	err := AddContact(dir, c2)
	if err == nil {
		t.Error("AddContact duplicate expected error")
	}
}

func TestGetContactNotFound(t *testing.T) {
	dir := make(Directory)
	_, err := GetContact(dir, "000-0000")
	if err == nil {
		t.Error("GetContact not found expected error")
	}
}

func TestUpdateContact(t *testing.T) {
	dir := make(Directory)
	c, _ := NewContact("Alice", "alice@example.com", "555-0101")
	AddContact(dir, c)
	updated, _ := NewContact("Alice Smith", "alice@new.com", "555-0101")
	err := UpdateContact(dir, "555-0101", updated)
	if err != nil {
		t.Fatalf("UpdateContact unexpected error: %v", err)
	}
	got, _ := GetContact(dir, "555-0101")
	if got.Name != "Alice Smith" || got.Email != "alice@new.com" {
		t.Errorf("UpdateContact = %+v; want {Alice Smith alice@new.com 555-0101}", got)
	}
}

func TestUpdateContactNotFound(t *testing.T) {
	dir := make(Directory)
	c, _ := NewContact("Alice", "a@x.com", "555-0101")
	err := UpdateContact(dir, "000-0000", c)
	if err == nil {
		t.Error("UpdateContact not found expected error")
	}
}

func TestDeleteContact(t *testing.T) {
	dir := make(Directory)
	c, _ := NewContact("Alice", "a@x.com", "555-0101")
	AddContact(dir, c)
	err := DeleteContact(dir, "555-0101")
	if err != nil {
		t.Fatalf("DeleteContact unexpected error: %v", err)
	}
	_, err = GetContact(dir, "555-0101")
	if err == nil {
		t.Error("GetContact after delete expected error")
	}
}

func TestDeleteContactNotFound(t *testing.T) {
	dir := make(Directory)
	err := DeleteContact(dir, "000-0000")
	if err == nil {
		t.Error("DeleteContact not found expected error")
	}
}

func TestListContactsEmpty(t *testing.T) {
	dir := make(Directory)
	list := ListContacts(dir)
	if list != nil {
		t.Errorf("ListContacts empty = %v; want nil", list)
	}
}

func TestListContactsOrder(t *testing.T) {
	dir := make(Directory)
	AddContact(dir, Contact{Name: "Zara", Email: "z@x.com", Phone: "555-0003"})
	AddContact(dir, Contact{Name: "Alice", Email: "a@x.com", Phone: "555-0001"})
	AddContact(dir, Contact{Name: "Bob", Email: "b@x.com", Phone: "555-0002"})
	list := ListContacts(dir)
	if len(list) != 3 {
		t.Fatalf("ListContacts length = %d; want 3", len(list))
	}
	if list[0].Name != "Alice" || list[1].Name != "Bob" || list[2].Name != "Zara" {
		t.Errorf("ListContacts order = %+v; want [Alice Bob Zara]", list)
	}
}

func TestSearchContacts(t *testing.T) {
	dir := make(Directory)
	AddContact(dir, Contact{Name: "Alice", Email: "alice@example.com", Phone: "555-0001"})
	AddContact(dir, Contact{Name: "Bob", Email: "bob@example.com", Phone: "555-0002"})

	results := SearchContacts(dir, "alice")
	if len(results) != 1 {
		t.Fatalf("SearchContacts('alice') length = %d; want 1", len(results))
	}
	if results[0].Name != "Alice" {
		t.Errorf("SearchContacts name = %q; want %q", results[0].Name, "Alice")
	}
}

func TestSearchContactsNoMatch(t *testing.T) {
	dir := make(Directory)
	AddContact(dir, Contact{Name: "Alice", Email: "a@x.com", Phone: "555-0001"})
	results := SearchContacts(dir, "nonexistent")
	if results != nil {
		t.Errorf("SearchContacts no match = %v; want nil", results)
	}
}
