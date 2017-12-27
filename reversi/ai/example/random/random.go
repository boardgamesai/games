package main

import (
	"github.com/boardgamesai/games/reversi"
	"github.com/boardgamesai/games/reversi/ai/driver"
	"github.com/boardgamesai/games/util"
)

type AI struct{}

func (ai *AI) GetMove(state driver.State) reversi.Move {
	moves := state.Board.PossibleMoves(state.Disc)
	return moves[util.RandInt(0, len(moves)-1)]
}
