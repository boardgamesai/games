package ulttictactoe

import (
	"math"
	"strings"

	"github.com/boardgamesai/games/tictactoe"
)

type Coords struct {
	Col int
	Row int
}

func getCoords(col, row int) *Coords {
	return &Coords{
		Col: col,
		Row: row,
	}
}

type Board struct {
	// SubGrids is shaped like this:
	// [0,2] [1,2] [2,2]
	// [0,1] [1,1] [2,1]
	// [0,0] [1,0] [2,0]
	SubGrids [3][3]*tictactoe.Board
	Grid     *tictactoe.Board // tracks the winners of the subgrids
	NextPlay *Coords          // location of a SubGrid; nil means can play anywhere
}

const Full = "-"

func NewBoard() *Board {
	subgrids := [3][3]*tictactoe.Board{}
	for i := range subgrids {
		for j := range subgrids[i] {
			subgrids[i][j] = &tictactoe.Board{}
		}
	}

	return &Board{
		SubGrids: subgrids,
		Grid:     &tictactoe.Board{},
	}
}

func (b *Board) IsValidMove(m Move) error {
	if m.Col > 2 || m.Col < 0 || m.Row > 2 || m.Row < 0 || m.SubCol > 2 || m.SubCol < 0 || m.SubRow > 2 || m.SubRow < 0 {
		return OutOfBoundsError{
			Move: m,
		}
	}

	if b.NextPlay != nil {
		// Are they obeying the next play rule?
		if m.Col != b.NextPlay.Col || m.Row != b.NextPlay.Row {
			return IllegalMoveError{
				Move: m,
			}
		}
	} else {
		// They can play anywhere, but are they trying to play in a subgrid that has a winner or is full?
		if b.Grid[m.Col][m.Row] != tictactoe.Empty {
			return IllegalMoveError{
				Move: m,
			}
		}
	}

	// Grid / Subgrid are good, but is it empty?
	if b.SubGrids[m.Col][m.Row][m.SubCol][m.SubRow] != tictactoe.Empty {
		return NotEmptyError{
			Move: m,
		}
	}

	return nil
}

// This returns winning moves for the subgrid, or notes that the subgrid is full with no winner
func (b *Board) ApplyMove(symbol string, m Move) ([]Move, bool, error) {
	winMoves := []Move{}
	subgridFilled := false

	if err := b.IsValidMove(m); err != nil {
		return winMoves, subgridFilled, err
	}

	subgrid := b.SubGrids[m.Col][m.Row]
	subgrid[m.SubCol][m.SubRow] = symbol

	if ok, moves := subgrid.HasWinner(); ok {
		// That move won that subgrid. Record it & send back the subgrid's winning moves
		b.Grid[m.Col][m.Row] = symbol
		for _, move := range moves {
			// Have to translate from tictactoe moves to ulttictactoe moves
			winMove := Move{
				Col:    m.Col,
				Row:    m.Row,
				SubCol: move.Col,
				SubRow: move.Row,
			}
			winMoves = append(winMoves, winMove)
		}
	}
	if subgrid.IsFull() {
		// That move filled that subgrid - only mark it if it's not already won
		if b.Grid[m.Col][m.Row] == tictactoe.Empty {
			b.Grid[m.Col][m.Row] = Full
		}
		subgridFilled = true
	}

	// Does the next play have a winner, or is it full? If so, they can play anywhere
	if b.Grid[m.SubCol][m.SubRow] != tictactoe.Empty {
		b.NextPlay = nil
	} else {
		b.NextPlay = getCoords(m.SubCol, m.SubRow)
	}

	return winMoves, subgridFilled, nil
}

// This returns winning moves for the tracking grid
func (b *Board) HasWinner() (bool, []tictactoe.Move) {
	b2 := b.Grid.DeepCopy()

	// This exists because we do something base tictactoe doesn't; we keep
	// track of "Full", but for win eval purposes we want "Full" to be ignored (e.g. "Empty")
	for i := range b2 {
		for j := range b2[i] {
			if b2[i][j] == Full {
				b2[i][j] = tictactoe.Empty
			}
		}
	}

	return b2.HasWinner()
}

func (b *Board) IsFull() bool {
	// When the tracking grid is all X / O / Full, then there's nowhere left to play
	return len(b.Grid.PossibleMoves()) == 0
}

func (b *Board) PossibleMoves() []Move {
	moves := []Move{}
	subgrids := map[*Coords]*tictactoe.Board{}

	// Figure out which subgrids can be played
	if b.NextPlay != nil {
		subgrids[b.NextPlay] = b.SubGrids[b.NextPlay.Col][b.NextPlay.Row]
	} else {
		for i := range b.Grid {
			for j := range b.Grid[i] {
				if b.Grid[i][j] == tictactoe.Empty {
					subgrids[getCoords(i, j)] = b.SubGrids[i][j]
				}
			}
		}
	}

	// Now grab all the empty spaces in the subgrid(s)
	for coords, subgrid := range subgrids {
		for i := range subgrid {
			for j := range subgrid[i] {
				if subgrid[i][j] == tictactoe.Empty {
					move := Move{
						Col:    coords.Col,
						Row:    coords.Row,
						SubCol: i,
						SubRow: j,
					}
					moves = append(moves, move)
				}
			}
		}
	}

	return moves
}

func (b *Board) String() string {
	rows := make([]string, 9)

	for y := 2; y >= 0; y-- { // Grid y
		for sy := 2; sy >= 0; sy-- { // Subgrid y
			line := make([]string, 3)
			for x := 0; x <= 2; x++ { // Grid x
				subrow := make([]string, 3)
				for sx := 0; sx <= 2; sx++ { // Subgrid x
					cell := b.SubGrids[x][y][sx][sy]
					if cell == tictactoe.Empty {
						cell = " "
					}
					subrow[sx] = cell
				}
				line[x] = strings.Join(subrow, ":")
			}
			rows[8-((y*3)+sy)] = strings.Join(line, "|")
		}
	}

	// Now add the horizontal grid lines
	s := ""
	for i, row := range rows {
		s += row + "\n"
		if i == 8 { // Don't print a line at the very bottom
			break
		}

		if (i+1)%3 == 0 {
			s += strings.Repeat("=", 17)
		} else {
			s += strings.Repeat("-", 17)
		}
		s += "\n"
	}

	// Also throw in the tracking grid
	s += "\n" + b.Grid.String()

	return s
}

func GetBoardFromString(s string) *Board {
	b := NewBoard()

	// This loads up all the subgrids
	for i, row := range strings.Split(s, "\n") {
		y := int(math.Floor((8 - float64(i)) / 3))
		sy := 2 - (i % 3)
		for x := 0; x <= 2; x++ {
			for sx := 0; sx <= 2; sx++ {
				cell := string(row[(x*3)+sx])
				if cell != "-" {
					b.SubGrids[x][y][sx][sy] = cell
				}
			}
		}
	}

	// Let's be extra clever and traverse the subgrids to set up the tracking grid
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			grid := b.SubGrids[i][j]
			if hasWinner, winMoves := grid.HasWinner(); hasWinner {
				// Have to dig who won it (X or O) out of winMoves
				move := winMoves[0]
				b.Grid[i][j] = grid[move.Col][move.Row]
			} else if grid.IsFull() {
				b.Grid[i][j] = Full
			}
		}
	}

	return b
}
