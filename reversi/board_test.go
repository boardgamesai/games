package reversi

import (
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func TestApplyMove(t *testing.T) {
	tests := []struct {
		board         string
		move          Move
		expectedBoard string
	}{
		{
			`
            ........
            ........
            ........
            ...WB...
            ...BW...
            ........
            ........
            ........
            `,
			Move{Col: 2, Row: 4},
			`
            ........
            ........
            ........
            ..BBB...
            ...BW...
            ........
            ........
            ........
            `,
		},
		{
			`
            ........
            ........
            ..W.....
            ..WBB...
            .WWBB...
            ..WBB...
            ...W....
            ........
            `,
			Move{Col: 1, Row: 2},
			`
            ........
            ........
            ..W.....
            ..WBB...
            .WBBB...
            .BBBB...
            ...W....
            ........
            `,
		},
		{
			`
            ....B...
            B...B...
            .BWWBB..
            .BWBBB..
            BWBBBB..
            W.BBBB..
            ..BBBB..
            ..BBBB..
            `,
			Move{Col: 0, Row: 1},
			`
            ....B...
            B...B...
            .BWWBB..
            .BWBBB..
            BWBBBB..
            B.BBBB..
            B.BBBB..
            ..BBBB..
            `,
		},
		{
			`
            ........
            ........
            ..BBBBB.
            ..BWWWB.
            ..BW.WB.
            ..BWWWB.
            ..BBBBB.
            ........
            `,
			Move{Col: 4, Row: 3},
			`
            ........
            ........
            ..BBBBB.
            ..BBBBB.
            ..BBBBB.
            ..BBBBB.
            ..BBBBB.
            ........
            `,
		},
	}

	for _, test := range tests {
		b := GetBoardFromString(trimBoard(test.board))
		_, err := b.ApplyMove(Black, test.move)
		if err != nil {
			t.Errorf("unexpected error applying move %s: %s", test.move, err)
			continue
		}

		b2 := GetBoardFromString(trimBoard(test.expectedBoard))
		if !boardsEqual(b, b2) {
			t.Errorf("boards not equal:\n%s\nand\n%s", b, b2)
		}
	}
}

func TestIsValidMove(t *testing.T) {
	tests := []struct {
		board    string
		col      int
		row      int
		expected error
	}{
		{
			`
            ........
            ........
            ........
            ...WB...
            ...BW...
            ........
            ........
            ........
            `,
			3,
			2,
			nil,
		},
		{
			`
            ........
            ........
            ........
            ...WB...
            ...BW...
            ........
            ........
            ........
            `,
			10,
			10,
			OutOfBoundsError{},
		},
		{
			`
            ........
            ........
            ........
            ...WB...
            ...BW...
            ........
            ........
            ........
            `,
			3,
			3,
			NotEmptyError{},
		},
		{
			`
            ........
            ........
            ........
            ...WB...
            ...BW...
            ........
            ........
            ........
            `,
			0,
			0,
			IllegalMoveError{},
		},
	}

	for _, test := range tests {
		board := GetBoardFromString(trimBoard(test.board))
		move := Move{
			Col: test.col,
			Row: test.row,
		}
		err := board.IsValidMove(White, move)
		if reflect.TypeOf(err) != reflect.TypeOf(test.expected) {
			t.Errorf("SetValidation board: %s set: [%d, %d] expected: %s got: %s", test.board, test.row, test.col, test.expected, err)
		}
	}
}

func TestPossibleMoves(t *testing.T) {
	tests := []struct {
		board         string
		expectedMoves []Move
	}{
		{
			`
            ........
            ........
            ........
            ...WB...
            ...BW...
            ........
            ........
            ........
            `,
			[]Move{{2, 4}, {3, 5}, {4, 2}, {5, 3}},
		},
		{
			`
            ........
            ........
            ..W.....
            ..WBB...
            .WWBB...
            ..WBB...
            ...W....
            ........
            `,
			[]Move{{0, 3}, {1, 1}, {1, 2}, {1, 4}, {1, 5}, {1, 6}, {2, 0}, {3, 0}},
		},
		{
			`
            ........
            ........
            ..BBBB..
            B.BWWW..
            BBWBBB..
            BBB.....
            ........
            ........
            `,
			[]Move{{6, 3}, {6, 4}, {6, 5}},
		},
		{
			`
            ....B...
            B...B...
            .BWWBB..
            .BWBBB..
            BWBBBB..
            W.BBBB..
            ..BBBB..
            ..BBBB..
            `,
			[]Move{{0, 1}, {0, 4}, {1, 2}, {1, 6}, {2, 6}, {3, 6}},
		},
		{
			`
            .BBBBBBB
            .WWWWW.B
            WWWWWWWB
            WWWWWWWB
            WWWWWWWB
            WWWWWWWB
            WWWWWWWB
            .WWWWW..
            `,
			[]Move{},
		},
	}

	for _, test := range tests {
		b := GetBoardFromString(trimBoard(test.board))
		moves := b.PossibleMoves(Black)
		if len(moves) != len(test.expectedMoves) {
			t.Errorf("got %d moves, expected %d for board:\n%s", len(moves), len(test.expectedMoves), b)
			continue
		}

		for i, m := range moves {
			if test.expectedMoves[i] != m {
				t.Errorf("got unexpected move %s (expected %s) for board:\n%s", m, test.expectedMoves[i], b)
			}
		}
	}
}

func TestBoardToFromString(t *testing.T) {
	str1 := `
        ........
        ........
        ..BBBB..
        B.BWWW..
        BBWBBB..
        BBB.....
        ........
        ........
    `
	str1 = trimBoard(str1)
	b := GetBoardFromString(str1)
	str2 := trimBoard(b.String())

	if str1 != str2 {
		t.Errorf("boards do not match:\n%s\nand\n%s\n", str1, str2)
	}
}

func trimBoard(s string) string {
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(strings.TrimSpace(s), "\n")
}

func boardsEqual(b1, b2 *Board) bool {
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if b1.Grid[i][j] != b2.Grid[i][j] {
				return false
			}
		}
	}

	return true
}
