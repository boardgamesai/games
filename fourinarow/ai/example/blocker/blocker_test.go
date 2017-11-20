package main

import (
	"regexp"
	"testing"

	"github.com/boardgamesai/games/fourinarow"
	"github.com/boardgamesai/games/fourinarow/ai/driver"
	"github.com/boardgamesai/games/game"
)

func TestGetMove(t *testing.T) {
	tests := []struct {
		boardStr string
		expected int
	}{
		{`0000000
          0000000
          0000000
          0000000
          0000000
          0011100`, 1},
		{`0000000
          0000000
          0000000
          0000000
          0000000
          0211100`, 5},
		{`0000000
          0000000
          0000000
          1100000
          1210000
          2121000`, 0},
		{`0000000
          0000000
          0000000
          0000011
          0000121
          0001212`, 6},
		{`0000000
          0000000
          0000000
          0000010
          0000010
          0000010`, 5},
		{`0100000
          0100010
          0100010
          0200010
          0100020
          0100010`, 5},
		{`0000000
          0000000
          0000000
          0000000
          0000000
          0110100`, 3},
		{`0000000
          0111000
          0212100
          0122200
          0211100
          0121200`, 4},
	}

	opponent := fourinarow.Player{
		Player: game.Player{
			Order: 1,
		},
	}
	state := driver.State{
		Order:    2,
		Opponent: &opponent,
	}

	re := regexp.MustCompile(`\s+`)

	ai := AI{}

	for _, test := range tests {
		boardStr := re.ReplaceAllString(test.boardStr, "|")
		state.Board = fourinarow.GetBoardFromString(boardStr)
		move := ai.GetMove(state)
		if move.Col != test.expected {
			t.Errorf("Blocker GetMove board: %s expected: %d got: %s", boardStr, test.expected, move)
		}
	}
}
