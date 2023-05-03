package main

import (
	"github.com/boardgamesai/games/game/elements/card"
	"github.com/boardgamesai/games/hearts"
	"github.com/boardgamesai/games/hearts/ai/driver"
	"github.com/boardgamesai/games/util"
)

type AI struct{}

func (ai *AI) GetPass(state driver.State) hearts.PassMove {
	move := hearts.PassMove{
		Cards: []card.Card{},
	}

	// Lots of different ways to do this. Let's make a slice containing [0..12],
	// shuffle it, and use the first three items as our hand indexes to pass.
	nums := make([]int, 13)
	for i := 0; i < 13; i++ {
		nums[i] = i
	}
	util.Shuffle(nums)

	for i := 0; i < 3; i++ {
		move.Cards = append(move.Cards, state.Hand[nums[i]])
	}

	return move
}

func (ai *AI) GetPlay(state driver.State) hearts.PlayMove {
	plays := state.Hand.PossiblePlays(state.Trick, state.TrickCount, state.HeartsBroken)
	return hearts.PlayMove{
		Card: plays[util.RandInt(0, len(plays)-1)],
	}
}
