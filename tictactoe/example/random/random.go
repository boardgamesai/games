package main

import (
	"github.com/boardgamesai/games/tictactoe"
	"github.com/boardgamesai/games/util"
)

func GetMove(state *tictactoe.State) tictactoe.Move {
	moves := state.Board.PossibleMoves()
	return moves[util.RandInt(0, len(moves)-1)]
}
