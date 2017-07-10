package tictactoe

import (
	"errors"
	"strings"
)

type Board struct {
	// board is shaped like this:
	// [0,0] [1,0] [2,0]
	// [0,1] [1,1] [2,1]
	// [0,2] [1,2] [2,2]
	Grid [3][3]string
}

const Empty = ""

var (
	ErrOutOfBounds = errors.New("out of bounds")
	ErrNotEmpty    = errors.New("not empty")
)

func (b *Board) IsValidMove(m Move) error {
	if m.Row > 2 || m.Row < 0 || m.Col > 2 || m.Col < 0 {
		return ErrOutOfBounds
	}

	if b.Grid[m.Col][m.Row] != Empty {
		return ErrNotEmpty
	}

	return nil
}

func (b *Board) HasWinner() bool {
	// Check rows & columns
	for i := 0; i < 3; i++ {
		row := [3]string{}
		col := [3]string{}
		for j := 0; j < 3; j++ {
			row[j] = b.Grid[j][i]
			col[j] = b.Grid[i][j]
		}

		if isThreeInARow(row) || isThreeInARow(col) {
			return true
		}
	}

	// Have to manually check diagonals
	diagonal1 := [3]string{b.Grid[0][0], b.Grid[1][1], b.Grid[2][2]}
	diagonal2 := [3]string{b.Grid[0][2], b.Grid[1][1], b.Grid[2][0]}
	if isThreeInARow(diagonal1) || isThreeInARow(diagonal2) {
		return true
	}

	return false
}

func isThreeInARow(vals [3]string) bool {
	return vals[0] != Empty && vals[0] == vals[1] && vals[1] == vals[2]
}

func (b *Board) IsFull() bool {
	return len(b.PossibleMoves()) == 0
}

func (b *Board) PossibleMoves() []Move {
	moves := []Move{}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if b.Grid[i][j] == Empty {
				moves = append(moves, Move{i, j})
			}
		}
	}

	return moves
}

func (b *Board) String() string {
	rows := make([]string, 3)

	for i := 0; i < 3; i++ {
		row := make([]string, 3)
		for j := 0; j < 3; j++ {
			cell := b.Grid[j][i]
			if cell == Empty {
				cell = " "
			}
			row[j] = cell
		}
		rows[i] = strings.Join(row, "|") + "\n"
	}

	return strings.Join(rows, "-----\n")
}

func GetStringFromBoard(b *Board) string {
	rows := make([]string, 3)

	for i := 0; i < 3; i++ {
		row := ""
		for j := 0; j < 3; j++ {
			cell := b.Grid[j][i]
			if cell == Empty {
				cell = " "
			}
			row += cell
		}
		rows[i] = row
	}

	return strings.Join(rows, "|")
}

func GetBoardFromString(s string) *Board {
	b := Board{}
	for i, row := range strings.Split(s, "|") {
		for j := 0; j < 3; j++ {
			cell := string(row[j])
			if cell != " " {
				b.Grid[j][i] = cell
			}
		}
	}

	return &b
}
