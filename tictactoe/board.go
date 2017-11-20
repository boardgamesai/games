package tictactoe

import (
	"errors"
	"strings"
)

type Board struct {
	// board is shaped like this:
	// [0,2] [1,2] [2,2]
	// [0,1] [1,1] [2,1]
	// [0,0] [1,0] [2,0]
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

func (b *Board) ApplyMove(symbol string, m Move) error {
	err := b.IsValidMove(m)
	if err != nil {
		return err
	}

	b.Grid[m.Col][m.Row] = symbol
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

func (b *Board) DeepCopy() *Board {
	newBoard := Board{}
	newBoard.Grid = b.Grid
	return &newBoard
}

// String is for nice visual display for humans
func (b *Board) String() string {
	rows := []string{}

	for i := 2; i >= 0; i-- {
		row := make([]string, 3)
		for j := 0; j <= 2; j++ {
			cell := b.Grid[j][i]
			if cell == Empty {
				cell = " "
			}
			row[j] = cell
		}
		rows = append(rows, strings.Join(row, "|")+"\n")
	}

	return strings.Join(rows, "-----\n")
}

// GetStringFromBoard represents a board as a simple string for test succinctness.
// For example, "X  |   |   " is a board with only an X in the top left.
func GetStringFromBoard(b *Board) string {
	rows := []string{}

	for i := 2; i >= 0; i-- {
		row := ""
		for j := 0; j <= 2; j++ {
			cell := b.Grid[j][i]
			if cell == Empty {
				cell = " "
			}
			row += cell
		}
		rows = append(rows, row)
	}

	return strings.Join(rows, "|")
}

func GetBoardFromString(s string) *Board {
	b := Board{}
	for i, row := range strings.Split(s, "|") {
		for j := 0; j <= 2; j++ {
			cell := string(row[j])
			if cell != " " {
				b.Grid[j][2-i] = cell
			}
		}
	}

	return &b
}
