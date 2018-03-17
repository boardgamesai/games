package amazons

import (
	"errors"
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

type Board struct {
	Grid [10][10]SpaceType
}

var (
	ErrOutOfBounds  = errors.New("out of bounds")
	ErrInvalidFrom  = errors.New("invalid from")
	ErrInvalidTo    = errors.New("invalid to")
	ErrInvalidArrow = errors.New("invalid arrow")
)

func NewBoard() *Board {
	b := Board{}
	b.Grid[0][3] = White
	b.Grid[3][0] = White
	b.Grid[6][0] = White
	b.Grid[9][3] = White
	b.Grid[0][6] = Black
	b.Grid[3][9] = Black
	b.Grid[6][9] = Black
	b.Grid[9][6] = Black
	return &b
}

func (b *Board) IsValidMove(s SpaceType, m Move) error {
	if offBoard(m.From.Col, m.From.Row) || offBoard(m.To.Col, m.To.Row) || offBoard(m.Arrow.Col, m.Arrow.Row) {
		return ErrOutOfBounds
	}

	if b.Grid[m.From.Col][m.From.Row] != s {
		return ErrInvalidFrom
	}

	if !moveListContains(b.PossibleMoves(m.From), m.To) {
		return ErrInvalidTo
	}

	if !moveListContains(b.PossibleArrows(m.From, m.To), m.Arrow) {
		return ErrInvalidArrow
	}

	return nil
}

func (b *Board) ApplyMove(s SpaceType, m Move) error {
	err := b.IsValidMove(s, m)
	if err != nil {
		return err
	}

	b.Grid[m.From.Col][m.From.Row] = Empty
	b.Grid[m.To.Col][m.To.Row] = s
	b.Grid[m.Arrow.Col][m.Arrow.Row] = Arrow
	return nil
}

func (b *Board) CanMove(player SpaceType) bool {
	return len(b.MovableQueens(player)) > 0
}

func (b *Board) MovableQueens(st SpaceType) []Space {
	spaces := []Space{}

	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			if b.Grid[i][j] == st && len(b.queenMoves(i, j)) > 0 {
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
	b2.Grid[to.Col][to.Row] = b2.Grid[from.Col][from.Row]
	b2.Grid[from.Col][from.Row] = Empty
	return b2.queenMoves(to.Col, to.Row)
}

func (b *Board) DeepCopy() *Board {
	newBoard := Board{}
	newBoard.Grid = b.Grid
	return &newBoard
}

func (b *Board) String() string {
	str := ""

	for i := 9; i >= 0; i-- {
		for j := 0; j <= 9; j++ {
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
			if offBoard(c, r) || b.Grid[c][r] != Empty {
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
			b.Grid[j][9-i] = SpaceType(space)
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
