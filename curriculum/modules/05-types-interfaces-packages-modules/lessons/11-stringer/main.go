package main

import "fmt"

type Status int

const (
	StatusPending Status = iota
	StatusActive
	StatusInactive
)

func (s Status) String() string {
	switch s {
	case StatusPending:
		return "pending"
	case StatusActive:
		return "active"
	case StatusInactive:
		return "inactive"
	default:
		return fmt.Sprintf("Status(%d)", int(s))
	}
}

type Point struct {
	X, Y int
}

func (p Point) String() string {
	return fmt.Sprintf("(%d, %d)", p.X, p.Y)
}

func (p Point) GoString() string {
	return fmt.Sprintf("Point{X: %d, Y: %d}", p.X, p.Y)
}

type Suit int

const (
	Hearts Suit = iota
	Diamonds
	Clubs
	Spades
)

func (s Suit) String() string {
	switch s {
	case Hearts:
		return "Hearts"
	case Diamonds:
		return "Diamonds"
	case Clubs:
		return "Clubs"
	case Spades:
		return "Spades"
	default:
		return fmt.Sprintf("Suit(%d)", int(s))
	}
}

type Rank int

const (
	Ace Rank = iota + 1
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

func (r Rank) String() string {
	names := map[Rank]string{
		Ace: "Ace", Two: "Two", Three: "Three", Four: "Four",
		Five: "Five", Six: "Six", Seven: "Seven", Eight: "Eight",
		Nine: "Nine", Ten: "Ten", Jack: "Jack", Queen: "Queen", King: "King",
	}
	if name, ok := names[r]; ok {
		return name
	}
	return fmt.Sprintf("Rank(%d)", int(r))
}

type Card struct {
	Suit Suit
	Rank Rank
}

func (c Card) String() string {
	return fmt.Sprintf("%s of %s", c.Rank, c.Suit)
}

func (c Card) GoString() string {
	return fmt.Sprintf("Card{Suit: %d, Rank: %d}", c.Suit, c.Rank)
}

func main() {
	fmt.Println(StatusPending)
	fmt.Println(StatusActive)
	fmt.Println(Status(99))

	p := Point{X: 3, Y: 4}
	fmt.Println(p)
	fmt.Printf("%v\n", p)
	fmt.Printf("%#v\n", p)

	cards := []Card{
		{Suit: Spades, Rank: Ace},
		{Suit: Hearts, Rank: King},
		{Suit: Diamonds, Rank: Seven},
	}
	for _, c := range cards {
		fmt.Printf("  %v  (%#v)\n", c, c)
	}
}
