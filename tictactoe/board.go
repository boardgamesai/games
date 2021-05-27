package tictactoe

import (
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

func (b *Board) IsValidMove(m Move) error {
	if m.Row > 2 || m.Row < 0 || m.Col > 2 || m.Col < 0 {
		return OutOfBoundsError{
			Move: m,
		}
	}

	if b.Grid[m.Col][m.Row] != Empty {
		return NotEmptyError{
			Move: m,
		}
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

func moves(col1, row1, col2, row2, col3, row3 int) []Move {
	return []Move{
		{Col: col1, Row: row1},
		{Col: col2, Row: row2},
		{Col: col3, Row: row3},
	}
}

func (b *Board) HasWinner() (bool, []Move) {
	allCoords := [][]Move{
		moves(0, 0, 1, 0, 2, 0), // rows
		moves(0, 1, 1, 1, 2, 1),
		moves(0, 2, 1, 2, 2, 2),
		moves(0, 0, 0, 1, 0, 2), // cols
		moves(1, 0, 1, 1, 1, 2),
		moves(2, 0, 2, 1, 2, 2),
		moves(0, 0, 1, 1, 2, 2), // diagonals
		moves(0, 2, 1, 1, 2, 0),
	}

	for _, coords := range allCoords {
		vals := [3]string{
			b.Grid[coords[0].Col][coords[0].Row],
			b.Grid[coords[1].Col][coords[1].Row],
			b.Grid[coords[2].Col][coords[2].Row],
		}
		if isThreeInARow(vals) {
			return true, coords
		}
	}

	return false, nil
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
