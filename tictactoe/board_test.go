package tictactoe

import (
	"reflect"
	"testing"
)

func TestIsValidMove(t *testing.T) {
	tests := []struct {
		boardStr string
		col      int
		row      int
		expected error
	}{
		{"   |   |   ", 0, 0, nil},
		{"   |   |   ", 1, 0, nil},
		{"   |   |   ", 2, 0, nil},
		{"   |   |   ", 0, 1, nil},
		{"   |   |   ", 1, 1, nil},
		{"   |   |   ", 2, 1, nil},
		{"   |   |   ", 0, 2, nil},
		{"   |   |   ", 1, 2, nil},
		{"   |   |   ", 2, 2, nil},
		{" XX|XXX|XXX", 0, 2, nil},
		{"X X|XXX|XXX", 1, 2, nil},
		{"XX |XXX|XXX", 2, 2, nil},
		{"XXX| XX|XXX", 0, 1, nil},
		{"XXX|X X|XXX", 1, 1, nil},
		{"XXX|XX |XXX", 2, 1, nil},
		{"XXX|XXX| XX", 0, 0, nil},
		{"XXX|XXX|X X", 1, 0, nil},
		{"XXX|XXX|XX ", 2, 0, nil},
		{"   |   |   ", 3, 2, OutOfBoundsError{}},
		{"   |   |   ", 2, 3, OutOfBoundsError{}},
		{"   |   |   ", 3, 3, OutOfBoundsError{}},
		{"   |   |   ", -1, 0, OutOfBoundsError{}},
		{"   |   |   ", 0, -1, OutOfBoundsError{}},
		{"X  |   |   ", 0, 2, NotEmptyError{}},
		{" X |   |   ", 1, 2, NotEmptyError{}},
		{"  X|   |   ", 2, 2, NotEmptyError{}},
		{"   |X  |   ", 0, 1, NotEmptyError{}},
		{"   | X |   ", 1, 1, NotEmptyError{}},
		{"   |  X|   ", 2, 1, NotEmptyError{}},
		{"   |   |X  ", 0, 0, NotEmptyError{}},
		{"   |   | X ", 1, 0, NotEmptyError{}},
		{"   |   |  X", 2, 0, NotEmptyError{}},
	}

	for _, test := range tests {
		board := GetBoardFromString(test.boardStr)
		move := Move{
			Col: test.col,
			Row: test.row,
		}
		err := board.IsValidMove(move)
		if reflect.TypeOf(err) != reflect.TypeOf(test.expected) {
			t.Errorf("IsValidMove board: %s set: [%d, %d] expected: %s got: %s", test.boardStr, test.row, test.col, test.expected, err)
		}
	}
}

func TestIsThreeInARow(t *testing.T) {
	tests := []struct {
		input    [3]string
		expected bool
	}{
		{[3]string{"", "", ""}, false},
		{[3]string{"X", "", ""}, false},
		{[3]string{"", "X", ""}, false},
		{[3]string{"", "", "X"}, false},
		{[3]string{"X", "", "X"}, false},
		{[3]string{"X", "X", ""}, false},
		{[3]string{"", "X", "X"}, false},
		{[3]string{"X", "", "O"}, false},
		{[3]string{"X", "X", "O"}, false},
		{[3]string{"O", "X", "O"}, false},
		{[3]string{"X", "X", "X"}, true},
	}

	for _, test := range tests {
		result := isThreeInARow(test.input)
		if result != test.expected {
			t.Errorf("isThreeInARow input: %s expected: %t got: %t", test.input, test.expected, result)
		}
	}
}

func TestPossibleMoves(t *testing.T) {
	tests := []struct {
		boardStr      string
		expectedMoves []Move
	}{
		{
			"   |   |   ",
			[]Move{{0, 0}, {0, 1}, {0, 2}, {1, 0}, {1, 1}, {1, 2}, {2, 0}, {2, 1}, {2, 2}},
		},
		{
			"   | X |   ",
			[]Move{{0, 0}, {0, 1}, {0, 2}, {1, 0}, {1, 2}, {2, 0}, {2, 1}, {2, 2}},
		},
		{
			"O  | X |   ",
			[]Move{{0, 0}, {0, 1}, {1, 0}, {1, 2}, {2, 0}, {2, 1}, {2, 2}},
		},
		{
			"O  | X |X  ",
			[]Move{{0, 1}, {1, 0}, {1, 2}, {2, 0}, {2, 1}, {2, 2}},
		},
		{
			"O O| X |X  ",
			[]Move{{0, 1}, {1, 0}, {1, 2}, {2, 0}, {2, 1}},
		},
		{
			"OXO| X |X  ",
			[]Move{{0, 1}, {1, 0}, {2, 0}, {2, 1}},
		},
		{
			"OXO| X |XO ",
			[]Move{{0, 1}, {2, 0}, {2, 1}},
		},
		{
			"OXO| XX|XO ",
			[]Move{{0, 1}, {2, 0}},
		},
		{
			"OXO|OXX|XO ",
			[]Move{{2, 0}},
		},
		{
			"OXO|OXX|XOX",
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
		winMoves  []Move
		hasWinner bool
	}{
		{"   |   |   ", nil, false},
		{"XXX|   |   ", moves(0, 2, 1, 2, 2, 2), true},
		{"   |XXX|   ", moves(0, 1, 1, 1, 2, 1), true},
		{"   |   |XXX", moves(0, 0, 1, 0, 2, 0), true},
		{"X  |X  |X  ", moves(0, 0, 0, 1, 0, 2), true},
		{" X | X | X ", moves(1, 0, 1, 1, 1, 2), true},
		{"  X|  X|  X", moves(2, 0, 2, 1, 2, 2), true},
		{"X  | X |  X", moves(0, 2, 1, 1, 2, 0), true},
		{"  X| X |X  ", moves(0, 0, 1, 1, 2, 2), true},
		{"XX |   |   ", nil, false},
		{" XX|   |   ", nil, false},
		{"X X|   |   ", nil, false},
		{"X  |X  |   ", nil, false},
		{"   |X  |X  ", nil, false},
		{"X  |   |X  ", nil, false},
		{"X X| X |   ", nil, false},
		{"   | X |X X", nil, false},
		{"XOX|XOO|OXX", nil, false},
		{"XOX|XOO|OOX", moves(1, 0, 1, 1, 1, 2), true},
	}

	for _, test := range tests {
		board := GetBoardFromString(test.boardStr)
		hasWinner, winMoves := board.HasWinner()
		if hasWinner != test.hasWinner {
			t.Errorf("hasWinner board: %s expected: %t got: %t", test.boardStr, test.hasWinner, hasWinner)
		}
		if !reflect.DeepEqual(winMoves, test.winMoves) {
			t.Errorf("winMoves board: %s expected: %+v got: %+v", test.boardStr, test.winMoves, winMoves)
		}
	}
}

func TestBoardToFromString(t *testing.T) {
	boardStr1 := "X X| X |OO "
	board := GetBoardFromString(boardStr1)
	boardStr2 := GetStringFromBoard(board)

	if boardStr1 != boardStr2 {
		t.Errorf("boards do not match: %s and %s", boardStr1, boardStr2)
	}
}
