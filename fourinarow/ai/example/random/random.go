package main

import (
	"github.com/boardgamesai/games/fourinarow"
	"github.com/boardgamesai/games/fourinarow/ai/driver"
	"github.com/boardgamesai/games/util"
)

type AI struct{}

func (ai *AI) GetMove(state driver.State) fourinarow.Move {
	moves := state.Board.PossibleMoves()
	return moves[util.RandInt(0, len(moves)-1)]
}
