package main

import (
	"github.com/boardgamesai/games/tictactoe"
	"github.com/boardgamesai/games/tictactoe/ai/driver"
	"github.com/boardgamesai/games/util"
)

type AI struct{}

func (ai *AI) GetMove(state driver.State) tictactoe.Move {
	moves := state.Board.PossibleMoves()
	return moves[util.RandInt(0, len(moves)-1)]
}
