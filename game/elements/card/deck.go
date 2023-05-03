package card

import (
	"fmt"

	"github.com/boardgamesai/games/util"
)

// List of types we'll allow here
type CardTypes interface {
	Card // Standard deck of 52 card
}

// We do this in two steps so we can mandate our accepted types have a String()
type Cardable interface {
	CardTypes
	fmt.Stringer
}

type Deck[C Cardable] struct {
	cards   []C
	current int
}

func NewDeck[C Cardable](cards []C) *Deck[C] {
	deck := Deck[C]{
		cards: cards,
	}
	deck.Shuffle()
	return &deck
}

func NewStandardDeck() *Deck[Card] {
	cards := make([]Card, 52)

	i := 0
	for _, suit := range Suits {
		for _, rank := range Ranks {
			cards[i] = NewCard(suit, rank)
			i++
		}
	}

	return NewDeck(cards)
}

func (d *Deck[C]) Shuffle() {
	util.Shuffle(d.cards)
	d.current = 0
}

func (d *Deck[C]) Count() int {
	return len(d.cards) - d.current
}

func (d *Deck[C]) DealCard() C {
	if d.current == len(d.cards) {
		// Past the end of the deck
		return C{}
	}

	card := d.cards[d.current]
	d.current++
	return card
}

func (d *Deck[C]) String() string {
	return fmt.Sprintf("%s", d.cards)
}
