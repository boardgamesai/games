package hearts

import (
	"errors"
	"fmt"
	"sort"

	"github.com/boardgamesai/games/hearts/card"
)

type Hand []card.Card

// Add has to be on a pointer because we are reassigning the underlying slice.
func (h *Hand) Add(c card.Card) {
	*h = append(*h, c)
}

// Remove has to be on a pointer because we are reassigning the underlying slice.
func (h *Hand) Remove(c card.Card) bool {
	for i, handCard := range *h {
		if handCard == c {
			*h = append((*h)[:i], (*h)[i+1:]...)
			return true
		}
	}

	return false
}

func (h Hand) Contains(c card.Card) bool {
	for _, card := range h {
		if card == c {
			return true
		}
	}

	return false
}

func (h Hand) Sort() {
	sort.Slice(h, func(i, j int) bool { return h[i].Index() < h[j].Index() })
}

func (h Hand) IsValidPass(cards []card.Card) error {
	if len(cards) != 3 {
		return errors.New("didn't get 3 pass cards")
	}

	// Make sure each card is actually in their hand
	passCards := map[card.Card]bool{}
	for _, passCard := range cards {
		if !h.Contains(passCard) {
			return fmt.Errorf("passed card %s not in hand", passCard)
		}

		if passCards[passCard] {
			return fmt.Errorf("duplicated card %s", passCard)
		}
		passCards[passCard] = true
	}

	return nil
}

func (h Hand) PossiblePlays(trick []card.Card, trickCount int, heartsBroken bool) []card.Card {
	plays := []card.Card{}

	for _, c := range h {
		// If they have the two of clubs and this is the first trick, ignore all else, they must play this.
		if trickCount == 0 && c.Rank == card.Two && c.Suit == card.Clubs {
			return []card.Card{c}
		}

		if len(trick) > 0 {
			// What we can play is based on the first card of the trick.
			if c.Suit == trick[0].Suit {
				plays = append(plays, c)
			}
		} else {
			// We are leading the trick, so we're free to play anything (unless hearts haven't been broken yet).
			if c.Suit == card.Hearts && !heartsBroken {
				continue
			}
			plays = append(plays, c)
		}
	}

	// If we get here with empty plays, that means someone led with a suit we don't have.
	// We're free to play any card in our hand, unless it's the first trick, then no hearts/queen of spades.
	if len(plays) == 0 {
		for _, c := range h {
			if trickCount == 0 && (c.Suit == card.Hearts || (c.Rank == card.Queen && c.Suit == card.Spades)) {
				continue
			}
			plays = append(plays, c)
		}
	}

	// Tbis is to handle the EXTREMELY unlikely case where it's the first trick and we were dealt only
	// hearts/queen of spades. In this case they can play anything.
	if len(plays) == 0 {
		plays = h
	}

	return plays
}
