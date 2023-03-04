package ulttictactoe

import (
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func trimBoard(s string) string {
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(strings.TrimSpace(s), "\n")
}

func getMove(col, row, subCol, subRow int) Move {
	return Move{
		Col:    col,
		Row:    row,
		SubCol: subCol,
		SubRow: subRow,
	}
}

func TestIsValidMove(t *testing.T) {
	tests := []struct {
		nextPlay    *Coords
		move        Move
		expectedErr error
	}{
		{
			getCoords(0, 0),
			getMove(0, 0, 0, 0),
			nil,
		},
		{
			getCoords(0, 0),
			getMove(0, 3, 0, 0),
			OutOfBoundsError{},
		},
		{
			getCoords(0, 0),
			getMove(0, 0, 0, 3),
			OutOfBoundsError{},
		},
		{
			getCoords(0, 0),
			getMove(-1, 0, 0, 0),
			OutOfBoundsError{},
		},
		{
			getCoords(0, 0),
			getMove(0, 0, -1, 0),
			OutOfBoundsError{},
		},
		{
			getCoords(0, 0),
			getMove(1, 0, 0, 0),
			IllegalMoveError{},
		},
		{
			getCoords(0, 0),
			getMove(0, 1, 2, 0),
			IllegalMoveError{},
		},
		{
			nil,
			getMove(1, 1, 2, 0),
			IllegalMoveError{},
		},
		{
			nil,
			getMove(2, 0, 2, 0),
			IllegalMoveError{},
		},
		{
			nil,
			getMove(2, 2, 2, 2),
			NotEmptyError{},
		},
		{
			nil,
			getMove(1, 2, 1, 2),
			NotEmptyError{},
		},
	}

	for i, test := range tests {
		board := GetBoardFromString(trimBoard(`
			----OX--X
			---------
			X--------
			---XXO---
			---OXX---
			---XOO---
			--------O
			-------O-
			------O--
		`))
		board.NextPlay = test.nextPlay

		err := board.IsValidMove(test.move)
		if reflect.TypeOf(err) != reflect.TypeOf(test.expectedErr) {
			t.Errorf("IsValidMove[%d] nextPlay: %+v move: %+v expected: %s got: %s", i, test.nextPlay, test.move, reflect.TypeOf(test.expectedErr), reflect.TypeOf(err))
		}
	}
}

func TestApplyMove(t *testing.T) {
	tests := []struct {
		board            string
		move             Move
		winMoves         []Move
		subgridFilled    bool
		expectedBoard    string
		expectedNextPlay *Coords
	}{
		{
			`
				---------
				---------
				---------
				---------
				---------
				---------
				---------
				---------
				---------
			`,
			getMove(1, 1, 0, 2),
			[]Move{},
			false,
			`
				---------
				---------
				---------
				---X-----
				---------
				---------
				---------
				---------
				---------
			`,
			getCoords(0, 2),
		},
		{
			`
				XXO------
				OOX------
				XO-------
				---------
				---------
				---------
				---------
				---------
				---------
			`,
			getMove(0, 2, 2, 0),
			[]Move{},
			true,
			`
				XXO------
				OOX------
				XOX------
				---------
				---------
				---------
				---------
				---------
				---------
			`,
			getCoords(2, 0),
		},
		{
			`
				XXO------
				OOX------
				XX-------
				---------
				---------
				---------
				---------
				---------
				---------
			`,
			getMove(0, 2, 2, 0),
			[]Move{getMove(0, 2, 0, 0), getMove(0, 2, 1, 0), getMove(0, 2, 2, 0)},
			true,
			`
				XXO------
				OOX------
				XXX------
				---------
				---------
				---------
				---------
				---------
				---------
			`,
			getCoords(2, 0),
		},
		{
			`
				XXO------
				O-X------
				XX-------
				---------
				---------
				---------
				---------
				---------
				---------
			`,
			getMove(0, 2, 2, 0),
			[]Move{getMove(0, 2, 0, 0), getMove(0, 2, 1, 0), getMove(0, 2, 2, 0)},
			false,
			`
				XXO------
				O-X------
				XXX------
				---------
				---------
				---------
				---------
				---------
				---------
			`,
			getCoords(2, 0),
		},
		{
			`
				--O------
				XOX------
				XXO------
				---------
				---------
				---------
				---------
				---------
				---------
			`,
			getMove(0, 2, 0, 2),
			[]Move{getMove(0, 2, 0, 0), getMove(0, 2, 0, 1), getMove(0, 2, 0, 2)},
			false,
			`
				X-O------
				XOX------
				XXO------
				---------
				---------
				---------
				---------
				---------
				---------
			`,
			nil,
		},
		{
			`
				OOO------
				---------
				---------
				---------
				---------
				---------
				---------
				---------
				---------
			`,
			getMove(2, 0, 0, 2),
			[]Move{},
			false,
			`
				OOO------
				---------
				---------
				---------
				---------
				---------
				------X--
				---------
				---------
			`,
			nil,
		},
	}

	for i, test := range tests {
		board := GetBoardFromString(trimBoard(test.board))
		winMoves, subgridFilled, err := board.ApplyMove("X", test.move)
		if err != nil {
			t.Errorf("ApplyMove[%d] got err: %s", i, err)
		} else if !movesMatch(winMoves, test.winMoves) {
			t.Errorf("ApplyMove[%d] expected winMoves: %+v got: %+v", i, test.winMoves, winMoves)
		} else if subgridFilled != test.subgridFilled {
			t.Errorf("ApplyMove[%d] expected subgridFilled: %t got: %t", i, test.subgridFilled, subgridFilled)
		} else if !nextPlayMatches(board.NextPlay, test.expectedNextPlay) {
			t.Errorf("ApplyMove[%d] expected nextPlay: %+v got: %+v", i, test.expectedNextPlay, board.NextPlay)
		} else if board.String() != GetBoardFromString(trimBoard(test.expectedBoard)).String() {
			t.Errorf("ApplyMove[%d] expected board:\n\n%sgot:%s\n\n", i, test.expectedBoard, test.board)
		}
	}
}

func movesMatch(moves1, moves2 []Move) bool {
	if len(moves1) != len(moves2) {
		return false
	}

	moveMap := map[Move]bool{}
	for _, move := range moves1 {
		moveMap[move] = true
	}

	for _, move := range moves2 {
		if !moveMap[move] {
			return false
		}
	}

	return true
}

func nextPlayMatches(np1, np2 *Coords) bool {
	if np1 == nil && np2 == nil {
		return true
	} else if (np1 == nil && np2 != nil) || (np1 != nil && np2 == nil) {
		return false
	}

	return np1.Col == np2.Col && np1.Row == np2.Row
}

func TestPossibleMoves(t *testing.T) {
	tests := []struct {
		board    string
		nextPlay *Coords
		numMoves int
	}{
		{
			`
				---------
				---------
				---------
				---------
				---------
				---------
				---------
				---------
				---------
			`,
			nil,
			81,
		},
		{
			`
				---------
				---------
				---------
				---X-----
				---------
				---------
				---------
				---------
				---------
			`,
			getCoords(0, 2),
			9,
		},
		{
			`
				---------
				-O-------
				---------
				---X-----
				---------
				---------
				---------
				---------
				---------
			`,
			getCoords(1, 1),
			8,
		},
		{
			`
				---------
				O--------
				---------
				---X-----
				---------
				---------
				---------
				---------
				---------
			`,
			getCoords(0, 1),
			9,
		},
		{
			`
				O--------
				O--------
				O--------
				---------
				---------
				---------
				---------
				---------
				---------
			`,
			nil,
			72,
		},
	}

	for i, test := range tests {
		board := GetBoardFromString(trimBoard(test.board))
		board.NextPlay = test.nextPlay

		moves := board.PossibleMoves()
		if len(moves) != test.numMoves {
			t.Errorf("PossibleMoves[%d] expected %d moves, got %d", i, test.numMoves, len(moves))
		}
	}
}
