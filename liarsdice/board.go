package liarsdice

import (
	"fmt"

	"github.com/boardgamesai/games/game/elements/dice"
)

type ChallengeOutcome struct {
	Bid            DiceVal
	ActualQuantity int
	DiceChanges    map[*Player]int
}

type Board struct {
	Bid        DiceVal
	Quantity   int
	Bidder     *Player
	DiceHidden map[*Player]*Dice
	DiceShown  map[*Player][]DiceVal
	Outcome    *ChallengeOutcome
}

func NewBoard(players []*Player) *Board {
	hidden := map[*Player]*Dice{}
	shown := map[*Player][]DiceVal{}

	for _, p := range players {
		hidden[p] = &Dice{dice.New(5, diceVals)}
		hidden[p].Roll()
		shown[p] = []DiceVal{}
	}

	return &Board{
		DiceHidden: hidden,
		DiceShown:  shown,
	}
}

func (b *Board) AllDice() []DiceVal {
	d := []DiceVal{}
	for p := range b.DiceHidden {
		d = append(d, b.DiceForPlayer(p)...)
	}
	return d
}

func (b *Board) DiceForPlayer(p *Player) []DiceVal {
	return append(b.DiceHidden[p].Values, b.DiceShown[p]...)
}

func (b *Board) IsValidMove(m Move) error {
	if m.Challenge {
		if b.Quantity == 0 {
			// Can't challenge when no one has made an opening bid yet
			return IllegalMoveError{m}
		}

		// Nothing below here applies if it's a challenge
		return nil
	}

	if m.Bid <= 0 || m.Bid > 6 || m.Quantity <= 0 || m.Quantity > len(b.AllDice()) {
		return OutOfBoundsError{m}
	}

	if m.Bid == Star {
		if b.Bid == Star {
			if m.Quantity <= b.Quantity {
				return IllegalMoveError{m}
			}
		} else {
			if m.Quantity*2 <= b.Quantity {
				return IllegalMoveError{m}
			}
		}
	} else {
		if b.Bid == Star {
			if b.Quantity*2 > m.Quantity {
				return IllegalMoveError{m}
			}
		} else {
			if m.Quantity == b.Quantity && m.Bid <= b.Bid {
				return IllegalMoveError{m}
			} else if m.Bid == b.Bid && m.Quantity <= b.Quantity {
				return IllegalMoveError{m}
			}
		}
	}

	return nil
}

func (b *Board) IsValidShow(m Move, p *Player) error {
	// You must have at least one die left to roll
	if len(m.ShowDice) >= b.DiceHidden[p].Count() {
		return IllegalMoveError{m}
	}

	// If there are ShowDice, we ensure they're a subset of the player's non-shown dice
	// First, we need to make a separate copy of this player's dice
	diceCopy := []DiceVal{}
	diceCopy = append(diceCopy, b.DiceHidden[p].Values...)

	// Now go over the shown dice, and for each one, ensure it's in our copy, and remove it
	for _, d := range m.ShowDice {
		found := false
		for i, d2 := range diceCopy {
			if d == d2 {
				diceCopy = append(diceCopy[:i], diceCopy[i+1:]...)
				found = true
				break
			}
		}

		if !found {
			return IllegalMoveError{m}
		}
	}

	return nil
}

func (b *Board) ApplyMove(m Move, p *Player) error {
	if err := b.IsValidMove(m); err != nil {
		return err
	}

	if len(m.ShowDice) > 0 {
		if err := b.IsValidShow(m, p); err != nil {
			return err
		}
	}

	if !m.Challenge {
		b.applyBidMove(m, p)
	} else {
		b.applyChallengeMove(p)
	}
	return nil
}

func (b *Board) applyBidMove(m Move, p *Player) {
	// First update the core board state
	b.Bid = m.Bid
	b.Quantity = m.Quantity
	b.Bidder = p
	b.Outcome = nil // Reset this from what it was for the last challenge

	// Now handle any shown dice
	if len(m.ShowDice) > 0 {
		b.moveShownDice(m, p)

		// Re-roll the dice this player has left
		b.DiceHidden[p].Roll()
	}
}

func (b *Board) moveShownDice(m Move, p *Player) {
	for _, d := range m.ShowDice {
		for i, d2 := range b.DiceHidden[p].Values {
			if d == d2 {
				b.DiceShown[p] = append(b.DiceShown[p], d)
				b.DiceHidden[p].Remove(i)
				break
			}
		}
	}
}

func (b *Board) applyChallengeMove(p *Player) {
	// How many dice are there, actually, of the bid?
	actual := 0
	for _, d := range b.AllDice() {
		if b.Bid == d || d == Star {
			actual++
		}
	}

	// Get the dice changes
	diceChanges := b.calculateDiceChanges(actual, p)

	// Store our outcome
	b.Outcome = &ChallengeOutcome{
		Bid:            b.Bid,
		ActualQuantity: actual,
		DiceChanges:    diceChanges,
	}

	// Flip shown dice back to hidden and apply the above changes
	b.restoreShownDice()
	b.applyDiceChanges(diceChanges)

	// Reset the board
	b.Bid = 0
	b.Bidder = nil
	b.Quantity = 0

	// Re-roll everyone's dice
	for _, d := range b.DiceHidden {
		d.Roll()
	}
}

func (b *Board) calculateDiceChanges(actual int, challenger *Player) map[*Player]int {
	diceChanges := map[*Player]int{}

	if actual < b.Quantity {
		// Bid was too high, bidder loses the difference
		change := min(b.Quantity-actual, len(b.DiceForPlayer(b.Bidder))) // Can only lose as many dice they have
		diceChanges[b.Bidder] = change * -1
	} else if actual > b.Quantity {
		// Bid was too low, challenger loses the difference
		change := min(actual-b.Quantity, len(b.DiceForPlayer(challenger)))
		diceChanges[challenger] = change * -1
	} else {
		// Bid was exactly right, challenger gives a die to the bidder
		diceChanges[challenger] = -1
		diceChanges[b.Bidder] = 1
	}

	return diceChanges
}

func (b *Board) restoreShownDice() {
	for p, d := range b.DiceShown {
		for i := 0; i < len(d); i++ {
			b.DiceHidden[p].Add(Star) // Value doesn't matter, will be re-rolled
		}
		b.DiceShown[p] = []DiceVal{}
	}
}

func (b *Board) applyDiceChanges(changes map[*Player]int) {
	for p, change := range changes {
		if change > 0 {
			// They gain dice (value doesn't matter, will be re-rolled)
			for change > 0 {
				b.DiceHidden[p].Add(Star)
				change--
			}
		} else {
			// They lose dice
			for change < 0 {
				b.DiceHidden[p].Remove(0)
				change++
			}
		}
	}
}

func (b *Board) String() string {
	return fmt.Sprintf("bid: %s quantity: %d dice in play: %d", b.Bid, b.Quantity, len(b.AllDice()))
}

func min(n1, n2 int) int {
	if n1 < n2 {
		return n1
	}
	return n2
}
