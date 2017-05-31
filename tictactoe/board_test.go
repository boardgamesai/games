package tictactoe

import (
	"fmt"
	"testing"
)

func TestSetValidation(t *testing.T) {
	tests := []struct {
		boardStr string
		row      int
		col      int
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
		{" XX|XXX|XXX", 0, 0, nil},
		{"X X|XXX|XXX", 1, 0, nil},
		{"XX |XXX|XXX", 2, 0, nil},
		{"XXX| XX|XXX", 0, 1, nil},
		{"XXX|X X|XXX", 1, 1, nil},
		{"XXX|XX |XXX", 2, 1, nil},
		{"XXX|XXX| XX", 0, 2, nil},
		{"XXX|XXX|X X", 1, 2, nil},
		{"XXX|XXX|XX ", 2, 2, nil},
		{"   |   |   ", 3, 2, ErrOutOfBounds},
		{"   |   |   ", 2, 3, ErrOutOfBounds},
		{"   |   |   ", 3, 3, ErrOutOfBounds},
		{"   |   |   ", -1, 0, ErrOutOfBounds},
		{"   |   |   ", 0, -1, ErrOutOfBounds},
		{"X  |   |   ", 0, 0, ErrNotEmpty},
		{" X |   |   ", 1, 0, ErrNotEmpty},
		{"  X|   |   ", 2, 0, ErrNotEmpty},
		{"   |X  |   ", 0, 1, ErrNotEmpty},
		{"   | X |   ", 1, 1, ErrNotEmpty},
		{"   |  X|   ", 2, 1, ErrNotEmpty},
		{"   |   |X  ", 0, 2, ErrNotEmpty},
		{"   |   | X ", 1, 2, ErrNotEmpty},
		{"   |   |  X", 2, 2, ErrNotEmpty},
	}

	for _, test := range tests {
		board := GetBoardFromString(test.boardStr)
		err := board.Set(test.row, test.col, "X")
		if err != test.expected {
			fmt.Printf("%s", board)
			t.Errorf("SetValidation board: %s set: [%d, %d] expected: %s got: %s", test.boardStr, test.row, test.col, test.expected, err)
		}
	}
}

func TestHasWinner(t *testing.T) {
	tests := []struct {
		boardStr string
		expected bool
	}{
		{"   |   |   ", false},
		{"XXX|   |   ", true},
		{"   |XXX|   ", true},
		{"   |   |XXX", true},
		{"X  |X  |X  ", true},
		{" X | X | X ", true},
		{"  X|  X|  X", true},
		{"X  | X |  X", true},
		{"  X| X |X  ", true},
		{"XX |   |   ", false},
		{" XX|   |   ", false},
		{"X X|   |   ", false},
		{"X  |X  |   ", false},
		{"   |X  |X  ", false},
		{"X  |   |X  ", false},
		{"X X| X |   ", false},
		{"   | X |X X", false},
		{"XOX|XOO|OXX", false},
		{"XOX|XOO|OOX", true},
	}

	for _, test := range tests {
		board := GetBoardFromString(test.boardStr)
		result := board.HasWinner()
		if result != test.expected {
			t.Errorf("IsWinner board: %s expected: %t got: %t", test.boardStr, test.expected, result)
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
