package fourinarow

import (
	"reflect"
	"testing"
)

func TestIsValidMove(t *testing.T) {
	tests := []struct {
		boardStr string
		col      int
		expected error
	}{
		{"0000000|0000000|0000000|0000000|0000000|0000000", 0, nil},
		{"0000000|0000000|0000000|0000000|0000000|0000000", 1, nil},
		{"0000000|0000000|0000000|0000000|0000000|0000000", 2, nil},
		{"0000000|0000000|0000000|0000000|0000000|0000000", 3, nil},
		{"0000000|0000000|0000000|0000000|0000000|0000000", 4, nil},
		{"0000000|0000000|0000000|0000000|0000000|0000000", 5, nil},
		{"0000000|0000000|0000000|0000000|0000000|0000000", 6, nil},
		{"0111111|1111111|1111111|1111111|1111111|1111111", 0, nil},
		{"1011111|1111111|1111111|1111111|1111111|1111111", 1, nil},
		{"1101111|1111111|1111111|1111111|1111111|1111111", 2, nil},
		{"1110111|1111111|1111111|1111111|1111111|1111111", 3, nil},
		{"1111011|1111111|1111111|1111111|1111111|1111111", 4, nil},
		{"1111101|1111111|1111111|1111111|1111111|1111111", 5, nil},
		{"1111110|1111111|1111111|1111111|1111111|1111111", 6, nil},
		{"0000000|0000000|0000000|0000000|0000000|0000000", -1, OutOfBoundsError{}},
		{"0000000|0000000|0000000|0000000|0000000|0000000", 7, OutOfBoundsError{}},
		{"1000000|1111111|1111111|1111111|1111111|1111111", 0, ColumnFullError{}},
		{"0100000|1111111|1111111|1111111|1111111|1111111", 1, ColumnFullError{}},
		{"0010000|1111111|1111111|1111111|1111111|1111111", 2, ColumnFullError{}},
		{"0001000|1111111|1111111|1111111|1111111|1111111", 3, ColumnFullError{}},
		{"0000100|1111111|1111111|1111111|1111111|1111111", 4, ColumnFullError{}},
		{"0000010|1111111|1111111|1111111|1111111|1111111", 5, ColumnFullError{}},
		{"0000001|1111111|1111111|1111111|1111111|1111111", 6, ColumnFullError{}},
	}

	for _, test := range tests {
		board := GetBoardFromString(test.boardStr)
		move := Move{
			Col: test.col,
		}
		err := board.IsValidMove(move)
		if reflect.TypeOf(err) != reflect.TypeOf(test.expected) {
			t.Errorf("IsValidMove board: %s move: %d expected: %s got: %s", test.boardStr, test.col, test.expected, err)
		}
	}
}

func TestIsFourInARow(t *testing.T) {
	tests := []struct {
		c1       int
		c2       int
		c3       int
		c4       int
		expected bool
	}{
		{0, 0, 0, 0, false},
		{1, 0, 0, 0, false},
		{0, 1, 0, 0, false},
		{0, 0, 1, 0, false},
		{0, 0, 0, 1, false},
		{1, 1, 0, 0, false},
		{1, 0, 1, 0, false},
		{1, 0, 0, 1, false},
		{1, 1, 1, 0, false},
		{1, 1, 0, 1, false},
		{2, 1, 1, 1, false},
		{1, 2, 1, 1, false},
		{1, 1, 2, 1, false},
		{1, 1, 1, 2, false},
		{1, 1, 1, 1, true},
		{2, 2, 2, 2, true},
	}

	for _, test := range tests {
		result := isFourInARow(test.c1, test.c2, test.c3, test.c4)
		if result != test.expected {
			t.Errorf("isFourInARow input: [%d,%d,%d,%d] expected: %t got: %t", test.c1, test.c2, test.c3, test.c4, test.expected, result)
		}
	}
}

func TestPossibleMoves(t *testing.T) {
	tests := []struct {
		boardStr      string
		expectedMoves []Move
	}{
		{
			"0000000|0000000|0000000|0000000|0000000|0000000",
			[]Move{{0}, {1}, {2}, {3}, {4}, {5}, {6}},
		},
		{
			"0100000|0100000|0100000|0100000|0100000|0100000",
			[]Move{{0}, {2}, {3}, {4}, {5}, {6}},
		},
		{
			"0000000|1111111|1111111|1111111|1111111|1111111",
			[]Move{{0}, {1}, {2}, {3}, {4}, {5}, {6}},
		},
		{
			"0001000|1111111|1111111|1111111|1111111|1111111",
			[]Move{{0}, {1}, {2}, {4}, {5}, {6}},
		},
		{
			"0001001|1111111|1111111|1111111|1111111|1111111",
			[]Move{{0}, {1}, {2}, {4}, {5}},
		},
		{
			"0101001|1111111|1111111|1111111|1111111|1111111",
			[]Move{{0}, {2}, {4}, {5}},
		},
		{
			"0101011|1111111|1111111|1111111|1111111|1111111",
			[]Move{{0}, {2}, {4}},
		},
		{
			"0111011|1111111|1111111|1111111|1111111|1111111",
			[]Move{{0}, {4}},
		},
		{
			"1111011|1111111|1111111|1111111|1111111|1111111",
			[]Move{{4}},
		},
		{
			"1111111|1111111|1111111|1111111|1111111|1111111",
			[]Move{},
		},
	}

	for _, test := range tests {
		board := GetBoardFromString(test.boardStr)
		moves := board.PossibleMoves()
		if !reflect.DeepEqual(moves, test.expectedMoves) {
			t.Errorf("PossibleMoves board: %s expected: %s got: %s", test.boardStr, test.expectedMoves, moves)
		}
	}
}

func TestHasWinner(t *testing.T) {
	tests := []struct {
		boardStr  string
		hasWinner bool
		winCoords []Coords
	}{
		{"0000000|0000000|0000000|0000000|0000000|0000000", false, nil},
		{"0000000|0000000|0000000|0000000|0000000|1111000", true, coords(0, 0, 1, 0, 2, 0, 3, 0)},
		{"0000000|0000000|0000000|0000000|0000000|0001111", true, coords(3, 0, 4, 0, 5, 0, 6, 0)},
		{"0000000|0000000|1000000|1000000|1000000|1000000", true, coords(0, 0, 0, 1, 0, 2, 0, 3)},
		{"0000000|0000000|0000001|0000001|0000001|0000001", true, coords(6, 0, 6, 1, 6, 2, 6, 3)},
		{"1000000|1000000|1000000|1000000|0000000|0000000", true, coords(0, 2, 0, 3, 0, 4, 0, 5)},
		{"0000001|0000001|0000001|0000001|0000000|0000000", true, coords(6, 2, 6, 3, 6, 4, 6, 5)},
		{"1111000|0000000|0000000|0000000|0000000|0000000", true, coords(0, 5, 1, 5, 2, 5, 3, 5)},
		{"0001111|0000000|0000000|0000000|0000000|0000000", true, coords(3, 5, 4, 5, 5, 5, 6, 5)},
		{"0000000|0001000|0001000|0001000|0001000|0000000", true, coords(3, 1, 3, 2, 3, 3, 3, 4)},
		{"0000000|0000000|0000000|0000000|0011110|0000000", true, coords(2, 1, 3, 1, 4, 1, 5, 1)},
		{"1000000|0100000|0010000|0001000|0000000|0000000", true, coords(3, 2, 2, 3, 1, 4, 0, 5)},
		{"0000001|0000010|0000100|0001000|0000000|0000000", true, coords(3, 2, 4, 3, 5, 4, 6, 5)},
		{"0000000|0000000|0001000|0010000|0100000|1000000", true, coords(0, 0, 1, 1, 2, 2, 3, 3)},
		{"0000000|0000000|0001000|0000100|0000010|0000001", true, coords(6, 0, 5, 1, 4, 2, 3, 3)},
		{"0000000|0100000|0010000|0001000|0000100|0000000", true, coords(4, 1, 3, 2, 2, 3, 1, 4)},
		{"0000000|0000100|0001000|0010000|0100000|0000000", true, coords(1, 1, 2, 2, 3, 3, 4, 4)},
		{"0000000|0000000|0000000|0000000|0000000|1112000", false, nil},
		{"0000000|0000000|0000000|0000000|0000000|0001211", false, nil},
		{"0000000|0000000|2000000|1000000|1000000|1000000", false, nil},
		{"0000000|0000000|0000001|0000001|0000001|0000002", false, nil},
		{"1000000|1000000|1000000|2000000|0000000|0000000", false, nil},
		{"0000001|0000001|0000001|0000002|0000000|0000000", false, nil},
		{"2111000|0000000|0000000|0000000|0000000|0000000", false, nil},
		{"0002111|0000000|0000000|0000000|0000000|0000000", false, nil},
		{"0000000|0001000|0001000|0001000|0002000|0000000", false, nil},
		{"0000000|0000000|0000000|0000000|0011210|0000000", false, nil},
		{"1000000|0100000|0010000|0002000|0000000|0000000", false, nil},
		{"0000001|0000010|0000100|0002000|0000000|0000000", false, nil},
		{"0000000|0000000|0001000|0020000|0100000|1000000", false, nil},
		{"0000000|0000000|0002000|0000100|0000010|0000001", false, nil},
		{"0000000|0200000|0010000|0001000|0000100|0000000", false, nil},
		{"0000000|0000100|0001000|0010000|0200000|0000000", false, nil},
		{"0000000|0000000|0000000|0000000|0000000|0111000", false, nil},
		{"0000000|0000000|0000000|0000000|0000000|0000111", false, nil},
		{"0000000|0000000|1000000|1000000|1000000|0000000", false, nil},
		{"0000000|0000000|0000000|0000001|0000001|0000001", false, nil},
		{"1000000|1000000|1000000|0000000|0000000|0000000", false, nil},
		{"0000000|0000001|0000001|0000001|0000000|0000000", false, nil},
		{"1110000|0000000|0000000|0000000|0000000|0000000", false, nil},
		{"0001110|0000000|0000000|0000000|0000000|0000000", false, nil},
		{"0000000|0001000|0001000|0001000|0000000|0000000", false, nil},
		{"0000000|0000000|0000000|0000000|0001110|0000000", false, nil},
		{"1000000|0100000|0010000|0000000|0000000|0000000", false, nil},
		{"0000001|0000010|0000100|0000000|0000000|0000000", false, nil},
		{"0000000|0000000|0001000|0010000|0100000|0000000", false, nil},
		{"0000000|0000000|0000000|0000100|0000010|0000001", false, nil},
		{"0000000|0100000|0010000|0001000|0000000|0000000", false, nil},
		{"0000000|0000100|0001000|0010000|0000000|0000000", false, nil},
	}

	for _, test := range tests {
		board := GetBoardFromString(test.boardStr)
		hasWinner, winMoves := board.HasWinner()
		if hasWinner != test.hasWinner {
			t.Errorf("IsWinner board: %s expected: %t got: %t", test.boardStr, test.hasWinner, hasWinner)
		}
		if !reflect.DeepEqual(winMoves, test.winCoords) {
			t.Errorf("winMoves board: %s expected: %+v got: %+v", test.boardStr, test.winCoords, winMoves)
		}
	}
}

func TestApplyMoveComputesRow(t *testing.T) {
	tests := []struct {
		boardStr    string
		col         int
		expectedRow int
	}{
		{"0000000|0000000|0000000|0000000|0000000|0000000", 0, 0},
		{"0000000|0000000|0000000|0000000|0000000|1000000", 0, 1},
		{"0000000|0000000|0000000|0000000|0000000|1000000", 1, 0},
		{"0000000|0000000|0000000|0000000|0000000|0100000", 1, 1},
		{"0000000|0000000|0000000|0000000|0100000|0100000", 1, 2},
		{"0000000|0000000|0000000|0100000|0100000|0100000", 1, 3},
		{"0000000|0000000|0100000|0100000|0100000|0100000", 1, 4},
		{"0000000|0100000|0100000|0100000|0100000|0100000", 1, 5},
	}

	for _, test := range tests {
		board := GetBoardFromString(test.boardStr)
		move := Move{
			Col: test.col,
		}
		row, _ := board.ApplyMove(2, move)
		if row != test.expectedRow {
			t.Errorf("TestApplyMoveComputesRow board: %s expected: %d got: %d", test.boardStr, test.expectedRow, row)
		}
	}
}

func TestBoardToFromString(t *testing.T) {
	boardStr1 := "0010000|0012000|0021000|1012100|2012200|1211212"
	board := GetBoardFromString(boardStr1)
	boardStr2 := GetStringFromBoard(board)

	if boardStr1 != boardStr2 {
		t.Errorf("boards do not match: %s and %s", boardStr1, boardStr2)
	}
}
