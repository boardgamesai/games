package main

import (
	"math"

	"github.com/boardgamesai/games/liarsdice"
	"github.com/boardgamesai/games/liarsdice/ai/driver"
	"github.com/boardgamesai/games/util"
)

type AI struct{}

func (ai *AI) GetMove(state driver.State) liarsdice.Move {
	m := liarsdice.Move{
		ShowDice: []liarsdice.DiceVal{},
	}

	if state.Bid == 0 {
		m.Bid = liarsdice.DiceVal(util.RandInt(1, 6))
		m.Quantity = 1
		return m
	}

	n := util.RandInt(1, 10)
	// log.Printf("n: %d", n)
	if n < 5 {
		// Increase the bid
		if n == 1 {
			m.Bid = liarsdice.Star
			if state.Bid == liarsdice.Star {
				m.Quantity = state.Quantity + 1
			} else {
				m.Quantity = int(math.Ceil(float64(state.Quantity+1) / 2))
			}
		} else {
			if state.Bid == 6 {
				m.Bid = liarsdice.DiceVal(util.RandInt(2, 6))
				m.Quantity = state.Quantity + 1
			} else if state.Bid == liarsdice.Star {
				m.Bid = liarsdice.DiceVal(util.RandInt(2, 6))
				m.Quantity = state.Quantity * 2
			} else {
				m.Bid = state.Bid + 1
				m.Quantity = state.Quantity
			}
		}

		if n < 3 && len(state.Dice) > 1 { // Must have 2+ dice to show any
			for _, d := range state.Dice {
				if d == m.Bid {
					m.ShowDice = append(m.ShowDice, d)
				}
			}
		}
	} else if n < 9 {
		// Increase the quantity
		m.Bid = state.Bid
		m.Quantity = state.Quantity + 1

		if n < 7 && len(state.Dice) > 1 { // Must have 2+ dice to show any
			for _, d := range state.Dice {
				if d == m.Bid {
					m.ShowDice = append(m.ShowDice, d)
				}
			}
		}
	} else {
		// Challenge
		m.Challenge = true
	}

	// If we've ended up with a quantity hihger than all dice left, switch to challenge
	count := 0
	for _, c := range state.DiceCounts {
		count += c
	}
	if count < m.Quantity {
		m.Challenge = true
	}

	return m
}
