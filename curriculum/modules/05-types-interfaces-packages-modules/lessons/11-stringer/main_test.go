package main

import "testing"

func TestStatusString(t *testing.T) {
	tests := []struct {
		s    Status
		want string
	}{
		{StatusPending, "pending"},
		{StatusActive, "active"},
		{StatusInactive, "inactive"},
	}
	for _, tc := range tests {
		got := tc.s.String()
		if got != tc.want {
			t.Errorf("Status(%d).String() = %q, want %q", int(tc.s), got, tc.want)
		}
	}
}

func TestStatusStringUnknown(t *testing.T) {
	s := Status(99)
	got := s.String()
	if got != "Status(99)" {
		t.Errorf("expected 'Status(99)', got %q", got)
	}
}

func TestPointString(t *testing.T) {
	p := Point{X: 3, Y: 4}
	if p.String() != "(3, 4)" {
		t.Errorf("expected '(3, 4)', got %q", p.String())
	}
}

func TestPointGoString(t *testing.T) {
	p := Point{X: 3, Y: 4}
	if p.GoString() != "Point{X: 3, Y: 4}" {
		t.Errorf("expected 'Point{X: 3, Y: 4}', got %q", p.GoString())
	}
}

func TestCardString(t *testing.T) {
	c := Card{Suit: Spades, Rank: Ace}
	if c.String() != "Ace of Spades" {
		t.Errorf("expected 'Ace of Spades', got %q", c.String())
	}
}

func TestCardGoString(t *testing.T) {
	c := Card{Suit: Hearts, Rank: King}
	expected := "Card{Suit: 0, Rank: 13}"
	if c.GoString() != expected {
		t.Errorf("expected %q, got %q", expected, c.GoString())
	}
}

func TestSuitString(t *testing.T) {
	tests := []struct {
		s    Suit
		want string
	}{
		{Hearts, "Hearts"},
		{Diamonds, "Diamonds"},
		{Clubs, "Clubs"},
		{Spades, "Spades"},
	}
	for _, tc := range tests {
		got := tc.s.String()
		if got != tc.want {
			t.Errorf("Suit(%d).String() = %q, want %q", int(tc.s), got, tc.want)
		}
	}
}

func TestRankString(t *testing.T) {
	tests := []struct {
		r    Rank
		want string
	}{
		{Ace, "Ace"},
		{Two, "Two"},
		{Jack, "Jack"},
		{Queen, "Queen"},
		{King, "King"},
	}
	for _, tc := range tests {
		got := tc.r.String()
		if got != tc.want {
			t.Errorf("Rank(%d).String() = %q, want %q", int(tc.r), got, tc.want)
		}
	}
}

func TestCardMultipleCards(t *testing.T) {
	cards := []Card{
		{Suit: Spades, Rank: Ace},
		{Suit: Hearts, Rank: King},
		{Suit: Diamonds, Rank: Seven},
	}
	expected := []string{
		"Ace of Spades",
		"King of Hearts",
		"Seven of Diamonds",
	}
	for i, c := range cards {
		if c.String() != expected[i] {
			t.Errorf("card %d: expected %q, got %q", i, expected[i], c.String())
		}
	}
}
