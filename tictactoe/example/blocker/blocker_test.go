package main

import (
	"testing"

	"github.com/boardgamesai/games/tictactoe"
)

func TestGetMove(t *testing.T) {
	tests := []struct {
		boardStr    string
		expectedCol int
		expectedRow int
	}{
		{"XX |   |   ", 2, 2},
		{"X X|   |   ", 1, 2},
		{" XX|   |   ", 0, 2},
		{"   |XX |   ", 2, 1},
		{"   |X X|   ", 1, 1},
		{"   | XX|   ", 0, 1},
		{"   |   |XX ", 2, 0},
		{"   |   |X X", 1, 0},
		{"   |   | XX", 0, 0},
		{"X  |X  |   ", 0, 0},
		{"X  |   |X  ", 0, 1},
		{"   |X  |X  ", 0, 2},
		{" X | X |   ", 1, 0},
		{" X |   | X ", 1, 1},
		{"   | X | X ", 1, 2},
		{"  X|  X|   ", 2, 0},
		{"  X|   |  X", 2, 1},
		{"   |  X|  X", 2, 2},
		{"X  | X |   ", 2, 0},
		{"X  |   |  X", 1, 1},
		{"   | X |  X", 0, 2},
		{"  X| X |   ", 0, 0},
		{"  X|   |X  ", 1, 1},
		{"   | X |X  ", 2, 2},
	}

	opponent := tictactoe.Player{
		Symbol: "X",
	}
	state := tictactoe.State{
		Symbol:   "O",
		Opponent: &opponent,
	}

	for _, test := range tests {
		state.Board = tictactoe.GetBoardFromString(test.boardStr)
		move := GetMove(&state)
		if move.Col != test.expectedCol || move.Row != test.expectedRow {
			t.Errorf("Blocker GetMove board: %s expected: [%d,%d] got: %s", test.boardStr, test.expectedCol, test.expectedRow, move)
		}
	}
}
