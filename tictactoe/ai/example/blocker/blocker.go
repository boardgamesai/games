package main

import (
	"github.com/boardgamesai/games/tictactoe"
	"github.com/boardgamesai/games/tictactoe/ai/driver"
	"github.com/boardgamesai/games/util"
)

type AI struct{}

func (ai *AI) GetMove(state driver.State) tictactoe.Move {
	allMoves := state.Board.PossibleMoves()
	for _, move := range allMoves {
		// See if the opponent playing this move would be a win for them.
		// We test on a copy of the board so we don't taint it for future moves (there is no UnapplyMove).
		board := state.Board.DeepCopy()
		board.ApplyMove(state.Opponent, move) // Ignore error because we know it's a possible move
		if board.HasWinner() {
			// We must block this move!
			return move
		}
	}

	// If we make it here, there's nothing to block, so do a random move.
	return allMoves[util.RandInt(0, len(allMoves)-1)]
}
