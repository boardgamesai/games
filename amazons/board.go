package amazons

import (
	"fmt"
	"strings"
)

type Space struct {
	Col int
	Row int
}

func (s Space) String() string {
	return fmt.Sprintf("[%d,%d]", s.Col, s.Row)
}

func (s Space) Equals(s2 Space) bool {
	return s.Col == s2.Col && s.Row == s2.Row
}

type SpaceType string

const (
	White = SpaceType("W")
	Black = SpaceType("B")
	Arrow = SpaceType("*")
	Empty = SpaceType("")
)

type Board [10][10]SpaceType

func NewBoard() *Board {
	b := Board{}
	b[0][3] = White
	b[3][0] = White
	b[6][0] = White
	b[9][3] = White
	b[0][6] = Black
	b[3][9] = Black
	b[6][9] = Black
	b[9][6] = Black
	return &b
}

func (b *Board) IsValidMove(s SpaceType, m Move) error {
	if offBoard(m.From.Col, m.From.Row) || offBoard(m.To.Col, m.To.Row) || offBoard(m.Arrow.Col, m.Arrow.Row) {
		return OutOfBoundsError{
			Move: m,
		}
	}

	if b[m.From.Col][m.From.Row] != s {
		return InvalidFromError{
			Move: m,
		}
	}

	if !moveListContains(b.PossibleMoves(m.From), m.To) {
		return InvalidToError{
			Move: m,
		}
	}

	if !moveListContains(b.PossibleArrows(m.From, m.To), m.Arrow) {
		return InvalidArrowError{
			Move: m,
		}
	}

	return nil
}

func (b *Board) ApplyMove(s SpaceType, m Move) error {
	err := b.IsValidMove(s, m)
	if err != nil {
		return err
	}

	b[m.From.Col][m.From.Row] = Empty
	b[m.To.Col][m.To.Row] = s
	b[m.Arrow.Col][m.Arrow.Row] = Arrow
	return nil
}

func (b *Board) CanMove(player SpaceType) bool {
	return len(b.MovableQueens(player)) > 0
}

func (b *Board) MovableQueens(st SpaceType) []Space {
	spaces := []Space{}

	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			if b[i][j] == st && len(b.queenMoves(i, j)) > 0 {
				spaces = append(spaces, Space{i, j})
			}
		}
	}

	return spaces
}

func (b *Board) PossibleMoves(from Space) []Space {
	return b.queenMoves(from.Col, from.Row)
}

func (b *Board) PossibleArrows(from, to Space) []Space {
	// Clone the board, apply the piece move, and get queen moves from there
	b2 := b.DeepCopy()
	b2[to.Col][to.Row] = b2[from.Col][from.Row]
	b2[from.Col][from.Row] = Empty
	return b2.queenMoves(to.Col, to.Row)
}

func (b *Board) DeepCopy() *Board {
	newBoard := Board{}
	newBoard = *b
	return &newBoard
}

func (b *Board) String() string {
	str := ""

	for i := 9; i >= 0; i-- {
		for j := 0; j <= 9; j++ {
			space := "."
			if b[j][i] != Empty {
				space = string(b[j][i])
			}
			str += space
		}
		str += "\n"
	}

	return str
}

func (b *Board) queenMoves(col, row int) []Space {
	spaces := []Space{}

	// Eight directions to check
	for i := 0; i < 8; i++ {
		c := col
		r := row
		cInc := 0
		rInc := 0

		for {
			switch i {
			case 0: // N
				rInc = 1
			case 1: // NE
				cInc = 1
				rInc = 1
			case 2: // E
				cInc = 1
			case 3: // SE
				cInc = 1
				rInc = -1
			case 4: // S
				rInc = -1
			case 5: // SW
				cInc = -1
				rInc = -1
			case 6: // W
				cInc = -1
			case 7: // NW
				cInc = -1
				rInc = 1
			}

			c += cInc
			r += rInc
			if offBoard(c, r) || b[c][r] != Empty {
				break
			}
			spaces = append(spaces, Space{c, r})
		}
	}

	return spaces
}

func GetBoardFromString(s string) *Board {
	b := Board{}

	for i, row := range strings.Split(s, "\n") {
		for j := 0; j < 10; j++ {
			space := string(row[j])
			if space == "." {
				space = ""
			}
			b[j][9-i] = SpaceType(space)
		}
	}

	return &b
}

func moveListContains(spaces []Space, space Space) bool {
	found := false
	for _, s := range spaces {
		if s.Equals(space) {
			found = true
			break
		}
	}
	return found
}

func offBoard(col, row int) bool {
	return col > 9 || col < 0 || row > 9 || row < 0
}
