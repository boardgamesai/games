package fourinarow

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Board is 7 wide by 6 high with [0,0] in the lower left and [6,5] in the top right
type Board struct {
	Grid [7][6]int
}

const Empty = 0

var (
	ErrOutOfBounds = errors.New("out of bounds")
	ErrColumnFull  = errors.New("column full")
)

func (b *Board) IsValidMove(m Move) error {
	if m.Col > 6 || m.Col < 0 {
		return ErrOutOfBounds
	}

	if b.Grid[m.Col][5] != Empty {
		return ErrColumnFull
	}

	return nil
}

func (b *Board) ApplyMove(p *Player, m Move) error {
	err := b.IsValidMove(m)
	if err != nil {
		return err
	}

	for i := 0; i < 6; i++ {
		if b.Grid[m.Col][i] == Empty {
			b.Grid[m.Col][i] = p.Order
			break
		}
	}

	return nil
}

func (b *Board) HasWinner() bool {
	// First check the columns
	for i := 0; i <= 6; i++ {
		for j := 0; j <= 2; j++ {
			if isFourInARow(b.Grid[i][j], b.Grid[i][j+1], b.Grid[i][j+2], b.Grid[i][j+3]) {
				return true
			}
		}
	}

	// Next check the rows
	for i := 0; i <= 5; i++ {
		for j := 0; j <= 3; j++ {
			if isFourInARow(b.Grid[j][i], b.Grid[j+1][i], b.Grid[j+2][i], b.Grid[j+3][i]) {
				return true
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
			if isFourInARow(b.Grid[i][j], b.Grid[i+1][j+1], b.Grid[i+2][j+2], b.Grid[i+3][j+3]) {
				return true
			}
		}
	}
	// Now the board on the right, the lower right start of each diagonal.
	for i := 3; i <= 6; i++ {
		for j := 0; j <= 2; j++ {
			if isFourInARow(b.Grid[i][j], b.Grid[i-1][j+1], b.Grid[i-2][j+2], b.Grid[i-3][j+3]) {
				return true
			}
		}
	}

	return false
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
		if b.Grid[i][5] == Empty {
			moves = append(moves, Move{i})
		}
	}

	return moves
}

func (b *Board) String() string {
	str := ""

	for i := 5; i >= 0; i-- {
		for j := 0; j <= 6; j++ {
			str += fmt.Sprintf("%d", b.Grid[j][i])
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
			row += fmt.Sprintf("%d", b.Grid[j][i])
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
			b.Grid[j][5-i] = cell
		}
	}

	return &b
}
