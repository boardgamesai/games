package main

import (
	"github.com/boardgamesai/games/ulttictactoe"
	"github.com/boardgamesai/games/ulttictactoe/ai/driver"
	"github.com/boardgamesai/games/util"
)

type AI struct{}

func (ai *AI) GetMove(state driver.State) ulttictactoe.Move {
	moves := state.Board.PossibleMoves()
	return moves[util.RandInt(0, len(moves)-1)]
}
