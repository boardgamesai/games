package tictactoe

import (
	"reflect"
	"testing"
)

func TestGetNewMovesForPlayer(t *testing.T) {
	tests := []struct {
		moves    []MoveLog
		expected []MoveLog
	}{
		{
			[]MoveLog{
				MoveLog{Move: Move{Col: 0, Row: 0}, Order: 1},
				MoveLog{Move: Move{Col: 1, Row: 0}, Order: 2},
			},
			[]MoveLog{
				MoveLog{Move: Move{Col: 1, Row: 0}, Order: 2},
			},
		},
		{
			[]MoveLog{
				MoveLog{Move: Move{Col: 0, Row: 0}, Order: 1},
				MoveLog{Move: Move{Col: 1, Row: 0}, Order: 2},
				MoveLog{Move: Move{Col: 2, Row: 0}, Order: 1},
				MoveLog{Move: Move{Col: 0, Row: 1}, Order: 2},
			},
			[]MoveLog{
				MoveLog{Move: Move{Col: 0, Row: 1}, Order: 2},
			},
		},
		{
			[]MoveLog{},
			[]MoveLog{},
		},
	}

	game := Game{}
	player := Player{
		Order: 1,
	}
	for _, test := range tests {
		game.Moves = test.moves
		newMoves := game.GetNewMovesForPlayer(&player)
		equal := reflect.DeepEqual(newMoves, test.expected)
		if !equal {
			t.Errorf("TestGetNewMovesForPlayer log: %+v expected: %+v got: %+v", test.moves, test.expected, newMoves)
		}
	}
}
