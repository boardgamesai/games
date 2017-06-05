package main

import (
	"github.com/boardgamesai/games/tictactoe"
	"github.com/boardgamesai/games/util"
)

func GetMove(symbol string, board *tictactoe.Board) tictactoe.Move {
	moves := board.PossibleMoves()
	return moves[util.RandInt(0, len(moves)-1)]
}
