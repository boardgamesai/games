package main

import (
	"github.com/boardgamesai/games/tictactoe"
	"github.com/boardgamesai/games/util"
)

func GetMove(symbol string, board *tictactoe.Board) tictactoe.Move {
	// Look for a potential win, and if we find one, block it.
	empty := tictactoe.Empty
	var opponent string
	if symbol == "X" {
		opponent = "O"
	} else {
		opponent = "X"
	}

	// First check the rows and columns.
	for i := 0; i < 3; i++ {
		row := [3]string{}
		col := [3]string{}
		for j := 0; j < 3; j++ {
			row[j] = board.Grid[j][i]
			col[j] = board.Grid[i][j]
		}

		if row[0] == opponent && row[1] == opponent && row[2] == empty {
			return tictactoe.Move{2, i}
		} else if row[0] == opponent && row[1] == empty && row[2] == opponent {
			return tictactoe.Move{1, i}
		} else if row[0] == empty && row[1] == opponent && row[2] == opponent {
			return tictactoe.Move{0, i}
		}

		if col[0] == opponent && col[1] == opponent && col[2] == empty {
			return tictactoe.Move{i, 2}
		} else if col[0] == opponent && col[1] == empty && col[2] == opponent {
			return tictactoe.Move{i, 1}
		} else if col[0] == empty && col[1] == opponent && col[2] == opponent {
			return tictactoe.Move{i, 0}
		}
	}

	// Next check the diagonals.
	diag := [3]string{board.Grid[0][0], board.Grid[1][1], board.Grid[2][2]}
	if diag[0] == opponent && diag[1] == opponent && diag[2] == empty {
		return tictactoe.Move{2, 2}
	} else if diag[0] == opponent && diag[1] == empty && diag[2] == opponent {
		return tictactoe.Move{1, 1}
	} else if diag[0] == empty && diag[1] == opponent && diag[2] == opponent {
		return tictactoe.Move{0, 0}
	}

	diag = [3]string{board.Grid[0][2], board.Grid[1][1], board.Grid[2][0]}
	if diag[0] == opponent && diag[1] == opponent && diag[2] == empty {
		return tictactoe.Move{2, 0}
	} else if diag[0] == opponent && diag[1] == empty && diag[2] == opponent {
		return tictactoe.Move{1, 1}
	} else if diag[0] == empty && diag[1] == opponent && diag[2] == opponent {
		return tictactoe.Move{0, 2}
	}

	// If we make it here, there's nothing to block, so do a random move.
	moves := board.PossibleMoves()
	return moves[util.RandInt(0, len(moves)-1)]
}
