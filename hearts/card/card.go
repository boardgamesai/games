package card

import "fmt"

type Suit string

const (
	Clubs    = Suit("C")
	Diamonds = Suit("D")
	Hearts   = Suit("H")
	Spades   = Suit("S")
)

type Rank string

const (
	Two   = Rank("2")
	Three = Rank("3")
	Four  = Rank("4")
	Five  = Rank("5")
	Six   = Rank("6")
	Seven = Rank("7")
	Eight = Rank("8")
	Nine  = Rank("9")
	Ten   = Rank("T")
	Jack  = Rank("J")
	Queen = Rank("Q")
	King  = Rank("K")
	Ace   = Rank("A")
)

type Card struct {
	Suit
	Rank
}

var (
	Suits   = []Suit{Clubs, Diamonds, Spades, Hearts}
	Ranks   = []Rank{Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King, Ace}
	suitMap = map[Suit]int{}
	rankMap = map[Rank]int{}
)

func init() {
	for i, suit := range Suits {
		suitMap[suit] = i
	}
	for i, rank := range Ranks {
		rankMap[rank] = i
	}
}

func New(suit Suit, rank Rank) Card {
	return Card{
		Suit: suit,
		Rank: rank,
	}
}

// Index is used to sort/compare cards
func (c Card) Index() int {
	return (suitMap[c.Suit] * 13) + rankMap[c.Rank]
}

func (c Card) String() string {
	return fmt.Sprintf("%s%s", c.Rank, c.Suit)
}

// FromString takes input like "4C" or "JH" and returns a Card
func FromString(s string) Card {
	return Card{
		Suit: Suit(s[1:]),
		Rank: Rank(s[0:1]),
	}
}
