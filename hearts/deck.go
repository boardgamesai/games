package hearts

import (
	"fmt"

	"github.com/boardgamesai/games/hearts/card"
	"github.com/boardgamesai/games/util"
)

type Deck struct {
	cards   []card.Card
	current int
}

func NewDeck() *Deck {
	deck := Deck{
		cards: make([]card.Card, 52),
	}
	deck.Shuffle()
	return &deck
}

func (d *Deck) Shuffle() {
	// Fill with all possible cards
	i := 0
	for _, suit := range card.Suits {
		for _, rank := range card.Ranks {
			d.cards[i] = card.New(suit, rank)
			i++
		}
	}

	util.Shuffle(d.cards)
	d.current = 0
}

func (d *Deck) Count() int {
	return 52 - d.current
}

func (d *Deck) DealCard() card.Card {
	if d.current == len(d.cards) {
		// Past the end of the deck
		return card.Card{}
	}

	card := d.cards[d.current]
	d.current++
	return card
}

func (d *Deck) String() string {
	return fmt.Sprintf("%s", d.cards)
}
