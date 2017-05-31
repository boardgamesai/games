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
	grid [3][3]string
}

const Empty = ""

var (
	ErrOutOfBounds = errors.New("out of bounds")
	ErrNotEmpty    = errors.New("not empty")
)

func (b *Board) Get(col, row int) string {
	return b.grid[col][row]
}

func (b *Board) Set(col, row int, symbol string) error {
	if row > 2 || row < 0 || col > 2 || col < 0 {
		return ErrOutOfBounds
	}

	if b.Get(col, row) != Empty {
		return ErrNotEmpty
	}

	b.grid[col][row] = symbol
	return nil
}

func (b *Board) HasWinner() bool {
	// Check rows & columns
	for i := 0; i < 3; i++ {
		row := [3]string{}
		col := [3]string{}
		for j := 0; j < 3; j++ {
			row[j] = b.Get(j, i)
			col[j] = b.Get(i, j)
		}

		if isThreeInARow(row) || isThreeInARow(col) {
			return true
		}
	}

	// Have to manually check diagonals
	diagonal1 := [3]string{b.Get(0, 0), b.Get(1, 1), b.Get(2, 2)}
	diagonal2 := [3]string{b.Get(0, 2), b.Get(1, 1), b.Get(2, 0)}
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
			if b.Get(i, j) == Empty {
				moves = append(moves, Move{i, j})
			}
		}
	}

	return moves
}

func (b *Board) DeepCopy() Board {
	newBoard := Board{}
	newBoard.grid = b.grid
	return newBoard
}

func (b *Board) String() string {
	rows := make([]string, 3)

	for i := 0; i < 3; i++ {
		row := make([]string, 3)
		for j := 0; j < 3; j++ {
			cell := b.Get(j, i)
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
			symbol := b.Get(j, i)
			if symbol == Empty {
				symbol = " "
			}
			row += symbol
		}
		rows[i] = row
	}

	return strings.Join(rows, "|")
}

func GetBoardFromString(s string) *Board {
	b := Board{}
	for i, row := range strings.Split(s, "|") {
		for j := 0; j < 3; j++ {
			symbol := string(row[j])
			if symbol != " " {
				b.Set(j, i, symbol)
			}
		}
	}

	return &b
}
