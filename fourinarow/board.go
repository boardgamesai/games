package fourinarow

import (
	"fmt"
	"strconv"
	"strings"
)

// Board is 7 wide by 6 high with [0,0] in the lower left and [6,5] in the top right
type Board [7][6]int

const Empty = 0

func (b *Board) IsValidMove(m Move) error {
	if m.Col > 6 || m.Col < 0 {
		return OutOfBoundsError{
			Move: m,
		}
	}

	if b[m.Col][5] != Empty {
		return ColumnFullError{
			Move: m,
		}
	}

	return nil
}

func (b *Board) ApplyMove(order int, m Move) (int, error) {
	row := -1

	err := b.IsValidMove(m)
	if err != nil {
		return row, err
	}

	for i := 0; i < 6; i++ {
		if b[m.Col][i] == Empty {
			b[m.Col][i] = order
			row = i
			break
		}
	}

	return row, nil
}

func coords(col1, row1, col2, row2, col3, row3, col4, row4 int) []Coords {
	return []Coords{
		{Col: col1, Row: row1},
		{Col: col2, Row: row2},
		{Col: col3, Row: row3},
		{Col: col4, Row: row4},
	}
}

func (b *Board) HasWinner() (bool, []Coords) {
	// First check the columns
	for i := 0; i <= 6; i++ {
		for j := 0; j <= 2; j++ {
			if isFourInARow(b[i][j], b[i][j+1], b[i][j+2], b[i][j+3]) {
				return true, coords(i, j, i, j+1, i, j+2, i, j+3)
			}
		}
	}

	// Next check the rows
	for i := 0; i <= 5; i++ {
		for j := 0; j <= 3; j++ {
			if isFourInARow(b[j][i], b[j+1][i], b[j+2][i], b[j+3][i]) {
				return true, coords(j, i, j+1, i, j+2, i, j+3, i)
			}
		}
	}

	// Finally, the diagonals. These are tricky - only certain places on the board can fit
	// four in a row diagonally, specifically:
	// ---xxxx     xxxx---
	// --xxxxx     xxxxx--
	// -xxxxxx     xxxxxx-
	// ooooxx- and -xxoooo
	// oooox--     --xoooo
	// oooo---     ---oooo

	// First check the board on the left, the lower left start of each diagonal.
	for i := 0; i <= 3; i++ {
		for j := 0; j <= 2; j++ {
			if isFourInARow(b[i][j], b[i+1][j+1], b[i+2][j+2], b[i+3][j+3]) {
				return true, coords(i, j, i+1, j+1, i+2, j+2, i+3, j+3)
			}
		}
	}
	// Now the board on the right, the lower right start of each diagonal.
	for i := 3; i <= 6; i++ {
		for j := 0; j <= 2; j++ {
			if isFourInARow(b[i][j], b[i-1][j+1], b[i-2][j+2], b[i-3][j+3]) {
				return true, coords(i, j, i-1, j+1, i-2, j+2, i-3, j+3)
			}
		}
	}

	return false, nil
}

func isFourInARow(g1, g2, g3, g4 int) bool {
	return (g1 != Empty) && (g1 == g2) && (g2 == g3) && (g3 == g4)
}

func (b *Board) IsFull() bool {
	return len(b.PossibleMoves()) == 0
}

func (b *Board) PossibleMoves() []Move {
	moves := []Move{}
	for i := 0; i < 7; i++ {
		if b[i][5] == Empty {
			moves = append(moves, Move{i})
		}
	}

	return moves
}

func (b *Board) DeepCopy() *Board {
	newBoard := Board{}
	newBoard = *b
	return &newBoard
}

func (b *Board) String() string {
	str := ""

	for i := 5; i >= 0; i-- {
		for j := 0; j <= 6; j++ {
			str += fmt.Sprintf("%d", b[j][i])
		}
		str += "\n"
	}

	return str
}

// GetStringFromBoard represents a board as a simple string for test succinctness.
// Example where only the first player has moved, and gone in the middle column:
// "0000000|0000000|0000000|0000000|0000000|0001000"
func GetStringFromBoard(b *Board) string {
	rows := []string{}

	for i := 5; i >= 0; i-- {
		row := ""
		for j := 0; j <= 6; j++ {
			row += fmt.Sprintf("%d", b[j][i])
		}
		rows = append(rows, row)
	}

	return strings.Join(rows, "|")
}

func GetBoardFromString(s string) *Board {
	b := Board{}

	for i, row := range strings.Split(s, "|") {
		for j := 0; j <= 6; j++ {
			cell, _ := strconv.Atoi(string(row[j]))
			b[j][5-i] = cell
		}
	}

	return &b
}
