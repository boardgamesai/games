package main

import (
	"github.com/boardgamesai/games/fourinarow"
	"github.com/boardgamesai/games/util"
)

func GetMove(state *fourinarow.State) fourinarow.Move {
	moves := state.Board.PossibleMoves()
	return moves[util.RandInt(0, len(moves)-1)]
}
