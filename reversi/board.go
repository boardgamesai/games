package reversi

import (
	"errors"
	"strings"
)

type Disc string

const (
	White = Disc("W")
	Black = Disc("B")
	Empty = Disc("")
)

type Board struct {
	Grid [8][8]Disc
}

var (
	ErrOutOfBounds = errors.New("out of bounds")
	ErrNotEmpty    = errors.New("not empty")
	ErrInvalid     = errors.New("invalid move")
)

func NewBoard() *Board {
	b := Board{}
	b.Grid[3][3] = Black
	b.Grid[3][4] = White
	b.Grid[4][3] = White
	b.Grid[4][4] = Black
	return &b
}

func (b *Board) IsValidMove(d Disc, m Move) error {
	if m.Row > 7 || m.Row < 0 || m.Col > 7 || m.Col < 0 {
		return ErrOutOfBounds
	}

	if b.Grid[m.Col][m.Row] != Empty {
		return ErrNotEmpty
	}

	found := false
	for _, move := range b.PossibleMoves(d) {
		if m == move {
			found = true
			break
		}
	}

	if !found {
		return ErrInvalid
	}

	return nil
}

func (b *Board) ApplyMove(d Disc, m Move) error {
	err := b.IsValidMove(d, m)
	if err != nil {
		return err
	}

	b.Grid[m.Col][m.Row] = d

	opponent := opposite(d)
	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			newCol := m.Col + x
			newRow := m.Row + y
			if offBoard(newCol, newRow) || b.Grid[newCol][newRow] != opponent {
				continue
			}

			if b.checkLine(newCol, newRow, x, y) {
				b.flipLine(newCol, newRow, x, y)
			}
		}
	}

	return nil
}

func (b *Board) PossibleMoves(d Disc) []Move {
	moves := []Move{}
	opponent := opposite(d)
	moveMap := map[Move]bool{}

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			// We only look at playing adjacent to our opponent
			if b.Grid[i][j] != opponent {
				continue
			}

			// Check the eight spaces around this opponent
			for x := -1; x <= 1; x++ {
				for y := -1; y <= 1; y++ {
					col := i + x
					row := j + y
					// If adjacent space is off the grid or not empty, it's not a possible move
					if offBoard(col, row) || b.Grid[col][row] != Empty {
						continue
					}

					move := Move{
						Col: col,
						Row: row,
					}
					if !moveMap[move] && b.checkLine(i, j, -1*x, -1*y) {
						moves = append(moves, move)
						moveMap[move] = true
					}
				}
			}
		}
	}

	return moves
}

func (b *Board) Score() map[Disc]int {
	scores := map[Disc]int{
		White: 0,
		Black: 0,
	}

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if b.Grid[i][j] != Empty {
				scores[b.Grid[i][j]]++
			}
		}
	}

	return scores
}

func (b *Board) IsFull() bool {
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if b.Grid[i][j] == Empty {
				return false
			}
		}
	}

	return true
}

func (b *Board) DeepCopy() *Board {
	newBoard := Board{}
	newBoard.Grid = b.Grid
	return &newBoard
}

func (b *Board) String() string {
	str := ""

	for i := 7; i >= 0; i-- {
		for j := 0; j <= 7; j++ {
			space := "."
			if b.Grid[j][i] != Empty {
				space = string(b.Grid[j][i])
			}
			str += space
		}
		str += "\n"
	}

	return str
}

func (b *Board) checkLine(col, row, deltaCol, deltaRow int) bool {
	stopDisc := opposite(b.Grid[col][row])

	for {
		// Continue in one direction, and if we hit the edge of the board or an
		// empty, it's no good. If we hit the stop disc, we're good.
		col += deltaCol
		row += deltaRow
		if offBoard(col, row) || b.Grid[col][row] == Empty {
			return false
		} else if b.Grid[col][row] == stopDisc {
			return true
		}
	}
}

func (b *Board) flipLine(col, row, deltaCol, deltaRow int) {
	stopDisc := opposite(b.Grid[col][row])

	for b.Grid[col][row] != stopDisc {
		b.Grid[col][row] = opposite(b.Grid[col][row])
		col += deltaCol
		row += deltaRow
	}
}

func GetBoardFromString(s string) *Board {
	b := Board{}

	for i, row := range strings.Split(s, "\n") {
		for j := 0; j < 8; j++ {
			cell := string(row[j])
			if cell == "." {
				cell = ""
			}
			b.Grid[j][7-i] = Disc(cell)
		}
	}

	return &b
}

func offBoard(col, row int) bool {
	return col > 7 || col < 0 || row > 7 || row < 0
}

func opposite(d Disc) Disc {
	if d == White {
		return Black
	}
	return White
}
