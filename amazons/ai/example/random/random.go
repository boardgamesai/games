package main

import (
	"github.com/boardgamesai/games/amazons"
	"github.com/boardgamesai/games/amazons/ai/driver"
	"github.com/boardgamesai/games/util"
)

type AI struct{}

func (ai *AI) GetMove(state driver.State) amazons.Move {
	queens := state.Board.MovableQueens(state.Color)
	queen := queens[util.RandInt(0, len(queens)-1)]

	moves := state.Board.PossibleMoves(queen)
	move := moves[util.RandInt(0, len(moves)-1)]

	arrows := state.Board.PossibleArrows(queen, move)
	arrow := arrows[util.RandInt(0, len(arrows)-1)]

	return amazons.Move{
		From:  queen,
		To:    move,
		Arrow: arrow,
	}
}
