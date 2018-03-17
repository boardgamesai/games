package amazons

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

func space(col, row int) Space {
	return Space{
		Col: col,
		Row: row,
	}
}

func trimBoard(s string) string {
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(strings.TrimSpace(s), "\n")
}

func TestIsValidMove(t *testing.T) {
	tests := []struct {
		move        Move
		expectedErr error
	}{
		{
			Move{
				From:  space(3, 0),
				To:    space(6, 3),
				Arrow: space(6, 8),
			},
			nil,
		},
		{
			Move{
				From:  space(-1, 0),
				To:    space(6, 3),
				Arrow: space(6, 8),
			},
			ErrOutOfBounds,
		},
		{
			Move{
				From:  space(3, 0),
				To:    space(6, 10),
				Arrow: space(6, 8),
			},
			ErrOutOfBounds,
		},
		{
			Move{
				From:  space(3, 0),
				To:    space(6, 3),
				Arrow: space(6, -5),
			},
			ErrOutOfBounds,
		},
		{
			Move{
				From:  space(2, 0),
				To:    space(6, 3),
				Arrow: space(6, 8),
			},
			ErrInvalidFrom,
		},
		{
			Move{
				From:  space(0, 6),
				To:    space(6, 3),
				Arrow: space(6, 8),
			},
			ErrInvalidFrom,
		},
		{
			Move{
				From:  space(3, 0),
				To:    space(0, 1),
				Arrow: space(6, 8),
			},
			ErrInvalidTo,
		},
		{
			Move{
				From:  space(3, 0),
				To:    space(3, 0),
				Arrow: space(6, 8),
			},
			ErrInvalidTo,
		},
		{
			Move{
				From:  space(3, 0),
				To:    space(3, 9),
				Arrow: space(6, 8),
			},
			ErrInvalidTo,
		},
		{
			Move{
				From:  space(3, 0),
				To:    space(6, 3),
				Arrow: space(6, 9),
			},
			ErrInvalidArrow,
		},
		{
			Move{
				From:  space(3, 0),
				To:    space(6, 3),
				Arrow: space(6, 3),
			},
			ErrInvalidArrow,
		},
		{
			Move{
				From:  space(3, 0),
				To:    space(6, 3),
				Arrow: space(0, 4),
			},
			ErrInvalidArrow,
		},
	}

	board := GetBoardFromString(trimBoard(`
        ...B..B...
        ..........
        ..........
        B........B
        ..........
        ..........
        W........W
        ..........
        ..........
        ...W..W...
    `))

	for i, test := range tests {
		err := board.IsValidMove(White, test.move)
		if err != test.expectedErr {
			t.Errorf("(%d) expected error %s, got %s", i, test.expectedErr, err)
		}
	}
}

func TestMovableQueens(t *testing.T) {
	tests := []struct {
		board          string
		spaceType      SpaceType
		expectedSpaces []Space
	}{
		{
			`
            ...B..B...
            ..........
            ..........
            B........B
            ..........
            ..........
            W........W
            ..........
            ..........
            ...W..W...
            `,
			White,
			[]Space{
				space(0, 3),
				space(3, 0),
				space(6, 0),
				space(9, 3),
			},
		},
		{
			`
            ...B..B...
            ..........
            ..........
            B........B
            ..........
            ..........
            W........W
            ..........
            ..........
            ...W..W...
            `,
			Black,
			[]Space{
				space(0, 6),
				space(3, 9),
				space(6, 9),
				space(9, 6),
			},
		},
		{
			`
            ...B..B...
            ..........
            ..........
            B........B
            ..........
            ..........
            .*********
            .*W*W*W*W*
            .*********
            ..........
            `,
			White,
			[]Space{},
		},
		{
			`
            ...B..B...
            ..........
            ..........
            B........B
            ..........
            ..........
            .***.*****
            .*W*W*W*W*
            .*********
            ..........
            `,
			White,
			[]Space{
				space(4, 2),
			},
		},
		{
			`
            ...B..B...
            ..........
            ..........
            B........B
            ..........
            ..........
            .****.****
            .*W*W*W*W*
            .*********
            ..........
            `,
			White,
			[]Space{
				space(4, 2),
				space(6, 2),
			},
		},
		{
			`
            ..........
            ..........
            ..........
            ...BB*....
            ...BW*....
            ...*B*....
            ..........
            ..........
            ..........
            ..........
            `,
			White,
			[]Space{},
		},
		{
			`
            ..........
            ..........
            ..........
            ...BB*....
            ...BW*....
            ...*B*....
            ..........
            ..........
            ..........
            ..........
            `,
			Black,
			[]Space{
				space(3, 5),
				space(3, 6),
				space(4, 4),
				space(4, 6),
			},
		},
		{
			`
            W.*.......
            ***.......
            ..........
            ..........
            ..........
            ..........
            ..........
            ..........
            ..........
            ..........
            `,
			White,
			[]Space{
				space(0, 9),
			},
		},
	}

	for i, test := range tests {
		board := GetBoardFromString(trimBoard(test.board))
		spaces := board.MovableQueens(test.spaceType)

		err := compareSpaceSlices(spaces, test.expectedSpaces)
		if err != nil {
			t.Errorf("(%d) %s", i, err)
		}
	}
}

func TestPossibleArrows(t *testing.T) {
	tests := []struct {
		board          string
		from           Space
		to             Space
		expectedSpaces []Space
	}{
		{
			`
            ..........
            ..........
            ..........
            ..........
            ......W...
            ..*W*B.*..
            ..BW...*..
            ..B***.*..
            .....*....
            .....*....
            `,
			space(3, 3),
			space(6, 3),
			[]Space{
				space(6, 4),
				space(6, 2),
				space(6, 1),
				space(6, 0),
				space(5, 3),
				space(4, 3),
				space(3, 3),
			},
		},
		{
			`
            ..........
            ..........
            ..........
            ..........
            ......W...
            ..*W*B.*..
            ..BW...*..
            ..B***.*..
            .....*....
            .....*....
            `,
			space(3, 3),
			space(5, 3),
			[]Space{
				space(6, 4),
				space(7, 5),
				space(8, 6),
				space(9, 7),
				space(6, 3),
				space(6, 2),
				space(7, 1),
				space(8, 0),
				space(4, 3),
				space(3, 3),
			},
		},
		{
			`
            ..........
            ..........
            ..........
            ..........
            ......W...
            ..*W*B.*..
            ..BW...*..
            ..B***.*..
            .....*....
            .....*....
            `,
			space(3, 3),
			space(4, 3),
			[]Space{
				space(5, 3),
				space(6, 3),
				space(3, 3),
			},
		},
	}

	for i, test := range tests {
		board := GetBoardFromString(trimBoard(test.board))
		spaces := board.PossibleArrows(test.from, test.to)

		err := compareSpaceSlices(spaces, test.expectedSpaces)
		if err != nil {
			t.Errorf("(%d) %s", i, err)
		}
	}
}

func TestQueenMoves(t *testing.T) {
	tests := []struct {
		board          string
		space          Space
		expectedSpaces []Space
	}{
		{
			`
            ..........
            ..........
            .W........
            ..........
            ..........
            ..........
            ..........
            ..........
            ..........
            ..........
            `,
			space(1, 7),
			[]Space{
				space(1, 8),
				space(1, 9),
				space(2, 8),
				space(3, 9),
				space(2, 7),
				space(3, 7),
				space(4, 7),
				space(5, 7),
				space(6, 7),
				space(7, 7),
				space(8, 7),
				space(9, 7),
				space(2, 6),
				space(3, 5),
				space(4, 4),
				space(5, 3),
				space(6, 2),
				space(7, 1),
				space(8, 0),
				space(1, 6),
				space(1, 5),
				space(1, 4),
				space(1, 3),
				space(1, 2),
				space(1, 1),
				space(1, 0),
				space(0, 6),
				space(0, 7),
				space(0, 8),
			},
		},
		{
			`
            ..........
            ...*......
            ..........
            ....W.....
            ...WB.....
            ..*.*.....
            ...B......
            ..........
            ..........
            ..........
            `,
			space(3, 5),
			[]Space{
				space(3, 6),
				space(3, 7),
				space(3, 4),
				space(2, 5),
				space(1, 5),
				space(0, 5),
				space(2, 6),
				space(1, 7),
				space(0, 8),
			},
		},
		{
			`
            ..........
            ..........
            ..........
            ..........
            ..........
            *****.....
            *****.....
            **W**.....
            *****.....
            *****.....
            `,
			space(2, 2),
			[]Space{},
		},
		{
			`
            ..........
            ..........
            ..........
            ..........
            ..........
            *****.....
            *****.....
            **W**.....
            **.**.....
            *****.....
            `,
			space(2, 2),
			[]Space{space(2, 1)},
		},
		{
			`
            ..........
            ..........
            ..........
            ..........
            ..........
            *****.....
            *****.....
            **W**.....
            *.***.....
            *****.....
            `,
			space(2, 2),
			[]Space{space(1, 1)},
		},
		{
			`
            ..........
            ..........
            ..........
            ..........
            ..........
            *****.....
            *****.....
            *.W**.....
            *****.....
            *****.....
            `,
			space(2, 2),
			[]Space{space(1, 2)},
		},
		{
			`
            ..........
            ..........
            ..........
            ..........
            ..........
            *****.....
            *.***.....
            **W**.....
            *****.....
            *****.....
            `,
			space(2, 2),
			[]Space{space(1, 3)},
		},
		{
			`
            ..........
            ..........
            ..........
            ..........
            ..........
            *****.....
            **.**.....
            **W**.....
            *****.....
            *****.....
            `,
			space(2, 2),
			[]Space{space(2, 3)},
		},
		{
			`
            ..........
            ..........
            ..........
            ..........
            ..........
            *****.....
            ***.*.....
            **W**.....
            *****.....
            *****.....
            `,
			space(2, 2),
			[]Space{space(3, 3)},
		},
		{
			`
            ..........
            ..........
            ..........
            ..........
            ..........
            *****.....
            *****.....
            **W.*.....
            *****.....
            *****.....
            `,
			space(2, 2),
			[]Space{space(3, 2)},
		},
		{
			`
            ..........
            ..........
            ..........
            ..........
            ..........
            *****.....
            *****.....
            **W**.....
            ***.*.....
            *****.....
            `,
			space(2, 2),
			[]Space{space(3, 1)},
		},
	}

	for i, test := range tests {
		board := GetBoardFromString(trimBoard(test.board))
		spaces := board.queenMoves(test.space.Col, test.space.Row)

		err := compareSpaceSlices(spaces, test.expectedSpaces)
		if err != nil {
			t.Errorf("(%d) %s", i, err)
		}
	}
}

func TestBoardToFromString(t *testing.T) {
	str1 := `
        ..........
        ...*......
        ..........
        ....W.....
        ...WB.....
        ..*.*.....
        ...B......
        ..........
        ..........
        ..........
    `
	str1 = trimBoard(str1)
	b := GetBoardFromString(str1)
	str2 := trimBoard(b.String())

	if str1 != str2 {
		t.Errorf("boards do not match:\n%s\nand\n%s\n", str1, str2)
	}
}

func compareSpaceSlices(got, expected []Space) error {
	if len(got) != len(expected) {
		return fmt.Errorf("expected %d spaces (%s), got %d (%s)", len(expected), expected, len(got), got)
	}

	for j := 0; j < len(got); j++ {
		if !got[j].Equals(expected[j]) {
			return fmt.Errorf("expected space: %s got: %s", expected[j], got[j])
		}
	}

	return nil
}
